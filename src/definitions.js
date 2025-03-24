let keywords = [
    "package",
    "import"
];

let typeKeywords = [
    "MFunc"
];

class ASTNode {
    constructor(parent) {
        // A parent is definitely wanted.
        if(parent != undefined && parent != null && parent != 'root') {
            this.parent = parent;
            parent.appendChild(this);    
        } else {
            if(parent != null && parent != 'root') {
                console.error("ASTNode has no parent! If root node, specify with a parent of 'root' or null.");
            }
        }

        this.children = [];
    }

    appendChild(child) {
        this.children.push(child);
    }
}

class RootASTNode extends ASTNode {
    constructor() { super('root'); }
}

class Expr extends ASTNode {
    constructor(parent) {
        super(parent);
    }

    serialize() {
        
    }
}