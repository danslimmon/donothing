# The magic of 8

This procedure implements a little arithmetic trick involving some manipulation of the
user's phone number.

## Enter your phone number

Enter your phone number, without area code. Formatting doesn't matter.

**Outputs**:
    - `PhoneNumber` (string): Your phone number

## Multiply your phone number by 8

Treating your phone number as a single integer, multiply it by 8.

**Inputs**:
    - `PhoneNumber`

**Outputs**:
    - `PhoneNumberTimesEight` (string): Your phone number times 8

## Add up the digits

Add up all the digits in both numbers, and then add 8 to the result. If the resulting sum
has more than one digit, take that sum and add up _its_ digits. Repeat until there's a single
digit left. That digit should be 8.

**Inputs**:
    - `PhoneNumber`
    - `PhoneNumberTimesEight`

