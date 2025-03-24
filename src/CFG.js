/*
    What is a CFG?
    
        A CFG is a Context-Free-Grammar.
    
    How is it implemented?

        Every Expr has a value that can come next, and if it doesn't 
            meet one of those values that can come next, there should be an error thrown.
            When that next value is correct, you add it to the AST (Abstract Syntax Tree), 
            specifying for the later code-generation step. But in order to get to that, we 
            must first parse the CFG rules to allow them to be understood in a 
            "What comes next?" context.
*/

let CFG = {
    "Expr" : ""
}

let CFG_parsed = null;

function getWhatNext() {

}