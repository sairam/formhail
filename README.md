## The Flow
### Regular User Flow
1. A user comes in, enters his email address in his own form
1. User submits the form to `https://domain.com/email@example.com`
1. When the email address is not associated with the Page/Domain
  1. An email is sent to confirm the user containing how to add optional UID
  1. When the user clicks the confirmation link
    1. Domain/user combo is confirmed
    1. Redirects to enter additional notifications like slack/webhook
  1. When the user clicks on Report Spam link, email id is added to blacklist
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
* submission to `https://domain.com/uid`
  * configuration is done through login via email of link which is valid for one time use
  * form ID generation will get a name along with other options/notification settings
  * 6-char prefix is limited to the FQDN+domain
  * rate limits apply on FQDN+domain and email id separately in Basic/Free plan
  * any option from the POSTed form gets preference
* Slack & Webhook notifications
* Limit on notifications based on NotificationType+Creds(limit emails to 1k submissions per duration)
* Site Wide confirmation for multiple forms

### Templates
* Confirmation Email
* Submission Email

## Love to Have Features
1. An intermediate page will show with a captcha (optional)
1. Option to send one email is sent per hour/day to avoid spam with delimiters like horizontal line or 20 at a time
1. Custom Validation based on regexp like phone / email address
1. Option to mark a message as Spam will consider the IP address and subject
1. Auto Confirm message to requester
1. GYOC - get your own credentials for signed stuff like `from` and `to`
1. Support SMS notification
1. Support Attachments via FileStack / s3
1. Paid feature for s3 upload + bandwidth with option of signed link to download

## Deployment
* Set environment variables `set -a ; . .env ; set +a`
