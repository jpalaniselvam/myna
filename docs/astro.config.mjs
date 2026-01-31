// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
	site: 'https://jpalaniselvam.github.io',
	base: '/myna',
	integrations: [
		starlight({
			title: 'Myna',
			social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/jpalaniselvam/myna' }],
			sidebar: [
				{
					label: 'Introduction',
					items: [
						{ label: 'Why myna?', slug: 'introduction/why' },
					],
				},
				{
					label: 'Getting Started',
					items: [
						{ label: 'Installation', slug: 'getting-started/installation' },
						{ label: 'Basic Usage', slug: 'getting-started/basic-usage' },
					],
				},
				{
					label: 'Guides',
					items: [
						{ label: 'Design Overview', slug: 'guides/design-overview' },
						{ label: 'Authentication', slug: 'guides/authentication' },
						{ label: 'Collections', slug: 'guides/collections' },
						{ label: 'Environments', slug: 'guides/environments' },
						{ label: 'Actions', slug: 'guides/actions' },
						{ label: 'Payload Handling', slug: 'guides/payloads' },
						{ label: 'Variable Resolution', slug: 'guides/variable-resolution' },
						{
							label: 'AWS Services',
							items: [
								{ label: 'Lambda', slug: 'guides/services/lambda' },
								{ label: 'SQS', slug: 'guides/services/sqs' },
								{ label: 'SNS', slug: 'guides/services/sns' },
								{ label: 'EventBridge', slug: 'guides/services/eventbridge' },
								{ label: 'EC2', slug: 'guides/services/ec2' },
								{ label: 'S3', slug: 'guides/services/s3' },
								{ label: 'Step Functions', slug: 'guides/services/step-functions' },
								{ label: 'DynamoDB', slug: 'guides/services/dynamodb' },
								{ label: 'RDS', slug: 'guides/services/rds' },
								{ label: 'SES', slug: 'guides/services/ses' }
							],
						},
					],
				},
				{
					label: 'Reference',
					autogenerate: { directory: 'reference' },
				},
			],
		}),
	],
});
