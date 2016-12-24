## The Flow
### Standard User Flow
1. A user comes in, enters his email address in his own form
1. User submits the form to `https://domain.com/email@example.com`
1. When the email address is not associated with the Page/Domain
  1. An intermediate page will show with a prompt to enter an optional slack hook + captcha
  1. An email is sent to confirm along with an option to enter the slack url
  1. The email link confirms for the Page, Slack Hook, form ID
1. The user is redirect back to the same page if \_next is not set. \_next is limited to the same proto+FQDN
1. Once the confirmation is made, an email is sent to the user with the submitted form
1. When a user's customer comes in and submits the form, he gets redirected to \_next page or referral page

### Support options from FrontEnd
* Form Options
  * form submission as POST to `https://domain.com/email@example.com`
  * `\_replyto/email`
  * `\_next` page to redirect to
  * `\_subject` - Email subject
  * `\_cc` - extra email addresses comma separated
  * `\_format` - `plain` default is `html`
  * `\_gotcha` - honeypot with hidden. if filled, will ignore the request
  * sending via AJAX via Accept header to `application/json`
* submission to `https://domain.com/6-char/uid`
  * configuration is done through login via email of link which is valid for one time use
  * form ID generation will get a name along with other options/notification settings
  * 6-char prefix is limited to the FQDN+domain
  * rate limits apply on FQDN+domain or 6char prefix for consistency
  * any option from the UI gets preference
* GYOC - get your own credentials for signed stuff like `from` and `to`
* Slack & Email notifications
* Limited notifications based on Type+Creds(limit emails to 1k submissions per duration)
* Site Wide confirmation for multiple pages/forms
* Option to send one email is sent per hour/day to avoid spam with delimiters like horizontal line or 20 at a time

### Templates
* Confirmation Email
* Submission Email

## Love to Have Features
1. Custom Validation based on regexp like phone / email address
1. Option to mark a message as Spam will consider the IP address and subject

### Based on BYOC
1. Support SMS notification with permalink
1. Support for PushOver/Push Notifications?
1. Support Attachments to link via s3
1. If attachment size > 10MB, put in s3 and provide a signed link
