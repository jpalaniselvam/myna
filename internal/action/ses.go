package action

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"

	"github.com/jpalaniselvam/myna/internal/parser"
	baseTypes "github.com/jpalaniselvam/myna/internal/types"
)

type SESAction struct {
	baseTypes.BaseAction
	SES SESConfig `toml:"ses"`
}

type SESConfig struct {
	Source       string   `toml:"source"`
	ToAddresses  []string `toml:"to_addresses"`
	CcAddresses  []string `toml:"cc_addresses"`
	BccAddresses []string `toml:"bcc_addresses"`
	Subject      string   `toml:"subject"`
	HtmlBody     string   `toml:"html_body"`
	TextBody     string   `toml:"text_body"`
	EmailAddress string   `toml:"email_address"`
}

func runSESAction(cfg aws.Config, content []byte, payload []byte) (interface{}, error) {
	var action SESAction
	if err := parser.Unmarshal(content, &action); err != nil {
		return nil, fmt.Errorf("failed to parse ses action: %w", err)
	}

	client := ses.NewFromConfig(cfg)
	ctx := context.TODO()

	switch action.Kind {
	case baseTypes.KindSESSendEmail:
		return runSESSendEmail(ctx, client, action, payload)
	case baseTypes.KindSESVerifyEmailIdentity:
		return runSESVerifyEmailIdentity(ctx, client, action)
	default:
		return nil, fmt.Errorf("unsupported ses kind: %s", action.Kind)
	}
}

func runSESSendEmail(ctx context.Context, client *ses.Client, action SESAction, payload []byte) (interface{}, error) {
	if action.SES.Source == "" {
		return nil, fmt.Errorf("source is required")
	}
	if len(action.SES.ToAddresses) == 0 {
		return nil, fmt.Errorf("to_addresses is required")
	}
	if action.SES.Subject == "" {
		return nil, fmt.Errorf("subject is required")
	}

	// If TextBody is empty and we have a payload, use it as TextBody
	if action.SES.TextBody == "" && len(payload) > 0 {
		action.SES.TextBody = string(payload)
	}

	if action.SES.HtmlBody == "" && action.SES.TextBody == "" {
		return nil, fmt.Errorf("at least one of html_body or text_body must be specified")
	}

	input := &ses.SendEmailInput{
		Source: aws.String(action.SES.Source),
		Destination: &types.Destination{
			ToAddresses:  action.SES.ToAddresses,
			CcAddresses:  action.SES.CcAddresses,
			BccAddresses: action.SES.BccAddresses,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(action.SES.Subject),
			},
			Body: &types.Body{},
		},
	}

	if action.SES.HtmlBody != "" {
		input.Message.Body.Html = &types.Content{
			Data: aws.String(action.SES.HtmlBody),
		}
	}
	if action.SES.TextBody != "" {
		input.Message.Body.Text = &types.Content{
			Data: aws.String(action.SES.TextBody),
		}
	}

	resp, err := client.SendEmail(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to send email: %w", err)
	}

	return map[string]string{
		"message_id": aws.ToString(resp.MessageId),
	}, nil
}

func runSESVerifyEmailIdentity(ctx context.Context, client *ses.Client, action SESAction) (interface{}, error) {
	if action.SES.EmailAddress == "" {
		return nil, fmt.Errorf("email_address is required")
	}

	_, err := client.VerifyEmailIdentity(ctx, &ses.VerifyEmailIdentityInput{
		EmailAddress: aws.String(action.SES.EmailAddress),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to verify email identity: %w", err)
	}

	return map[string]string{
		"status": "verification email sent",
	}, nil
}
