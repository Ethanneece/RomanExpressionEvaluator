# Roman Expression Evaluator

This is a project created to parse expressions written in a mix of 
    English and Roman numerals.

Given a file:
    this program will read each expression 
    track the variable that it was written in
    and output the value of the final variable of the expressions

If the program finds an error in the code:
    The program will stop and output the error message and 
    where the error occured. 

Roman literal output limited to [1 - 3999]

Examples of running the program avalible in run.sh

## Error Messages 

These are the messages that may appear if something goes awry while parsing the expressions
Will also print out where the error happend for easy fixing. 

### LexicalErrMsg

Error that occurs when token does not match pattern of any valid sequence of characters. 

Error output: Quid dicis? You offend Caesar with your sloppy lexical habits!
Error code: -10

### SyntaxErrMsg

Error that occurs when intended structure of the expression is not valid.

Error output: Quid dicis? True Romans would not understand your syntax!
Error code: -100

### DeclErrMsg

Error that occurs when a variable being read is one that has not been declared yet. 

Error output: Quid dicis? Failure to declare allegiance to Caesar!
Error code: -200

### ZeroErrMsg

Error that occurs when the number zero comes up in calculations. 

Error output: Quid dicis? Arab merchants haven't left for India yet!
Error code: -500

### NegErrMsg

Error that occurs when a negative number comes up in calculation

Error output: Quid dicis? Caesar demands positive thoughts!
Error code: -300

## Operators

PLUS, MINUS, TIMES, DIVIDE, POWER, MODULO, EST

# Build Project
- make
# Run Project
- ./run.sh
# Clean Project
- make clean
