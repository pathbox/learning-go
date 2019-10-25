One Time Passwords (OTPs) are an mechanism to improve security over passwords alone. When a Time-based OTP (TOTP) is stored on a user's phone, and combined with something the user knows (Password), you have an easy on-ramp to Multi-factor authentication without adding a dependency on a SMS provider. This Password and TOTP combination is used by many popular websites including Google, Github, Facebook, Salesforce and many others.

The otp library enables you to easily add TOTPs to your own application, increasing your user's security against mass-password breaches and malware.

Because TOTP is standardized and widely deployed, there are many mobile clients and software implementations.

Generating QR Code images for easy user enrollment.
Time-based One-time Password Algorithm (TOTP) (RFC 6238): Time based OTP, the most commonly used method.
HMAC-based One-time Password Algorithm (HOTP) (RFC 4226): Counter based OTP, which TOTP is based upon.
Generation and Validation of codes for either algorithm

https://github.com/pquerna/otp

User Enrollment
For an example of a working enrollment work flow, Github has documented theirs, but the basics are:

1.Generate new TOTP Key for a User. key,_ := totp.Generate(...).
2.Display the Key's Secret and QR-Code for the User. key.Secret() and key.Image(...).
3.Test that the user can successfully use their TOTP. totp.Validate(...).
4.Store TOTP Secret for the User in your backend. key.Secret()
5.Provide the user with "recovery codes". (See Recovery Codes bellow)

Code Generation
In either TOTP or HOTP cases, use the GenerateCode function and a counter or time.Time struct to generate a valid code compatible with most implementations.
For uncommon or custom settings, or to catch unlikely errors, use GenerateCodeCustom in either module.
Validation
1. Prompt and validate User's password as normal.
2. If the user has TOTP enabled, prompt for TOTP passcode.
3. Retrieve the User's TOTP Secret from your backend.
4. Validate the user's passcode. totp.Validate(...)