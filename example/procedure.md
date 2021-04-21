# The magic of 8

This procedure implements a little arithmetic trick involving some manipulation of the
user's phone number.

- [Enter your phone number](#0-enter-your-phone-number)
- [Multiply your phone number by 8](#1-multiply-your-phone-number-by-8)
- [Add up the digits](#2-add-up-the-digits)

## (0) Enter your phone number

`root.inputPhoneNumber`
•
[Up](#the-magic-of-8)

Enter your phone number, without area code. Formatting doesn't matter.

**Outputs**:

  - `PhoneNumber` (string): Your phone number

## (1) Multiply your phone number by 8

`root.multiplyPhoneNumber`
•
[Up](#the-magic-of-8)

Treating your phone number as a single integer, multiply it by 8.

**Inputs**:

  - `PhoneNumber`

**Outputs**:

  - `PhoneNumberTimesEight` (string): Your phone number times 8

## (2) Add up the digits

`root.addDigits`
•
[Up](#the-magic-of-8)

Add up all the digits in both numbers, and then add 8 to the result. If the resulting sum
has more than one digit, take that sum and add up _its_ digits. Repeat until there's a single
digit left. That digit should be 8.

**Inputs**:

  - `PhoneNumber`
  - `PhoneNumberTimesEight`
