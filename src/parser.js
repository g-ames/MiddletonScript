function parse(tokens) {
    let ast = new RootASTNode();

    for(var i = 0; i < tokens.length; i++) {
        let token = tokens[i];

        console.log(token);
    }

    return ast;
}