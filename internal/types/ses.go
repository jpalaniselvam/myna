package types

type SESAction struct {
	BaseAction
	SES SESConfig `toml:"ses" json:"ses"`
}

type SESConfig struct {
	Source       string   `toml:"source" json:"source"`
	ToAddresses  []string `toml:"to_addresses" json:"to_addresses"`
	CcAddresses  []string `toml:"cc_addresses,omitempty" json:"cc_addresses,omitempty"`
	BccAddresses []string `toml:"bcc_addresses,omitempty" json:"bcc_addresses,omitempty"`
	Subject      string   `toml:"subject" json:"subject"`
	HtmlBody     string   `toml:"html_body,omitempty" json:"html_body,omitempty"`
	TextBody     string   `toml:"text_body,omitempty" json:"text_body,omitempty"`
	EmailAddress string   `toml:"email_address,omitempty" json:"email_address,omitempty"`
}
