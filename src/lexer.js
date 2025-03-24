function isNumericCharacter(char) {
    return !isNaN(Number(char)) && typeof Number(char) === 'number';
}

function isAlphabeticCharacter(char) {
    return /^[A-Za-z]+$/.test(char);
}

function getCharType(c) {
    if(c.trim() != c || c == " ") {
        return "whitespace";
    } else if(isNumericCharacter(c) || isAlphabeticCharacter(c)) {
        return "alnum";
    } else {
        return "other";
    }
}

function lex(source) {
    let in_quote = false;
    let char_type = "";
    let last_char_type = "";
    let token = "";
    let tokens = [];
    let line = 0;
    let tokenStart = 0;

    function pushCatagorizedToken(token) {
        if(token[0] == '"' && token[token.length - 1] == '"') {
            tokens.push({value: token, type: "string", at: line});
        } else {
            let ct = {value: token, type: getCharType(token[0]), at: line};
            
            if(keywords.includes(ct.value)) {
                ct.type = "keyword";
            } else if(typeKeywords.includes(ct.value)) {
                ct.type = "type";
            }

            tokens.push(ct);
        }
    }

    function pushToken(also) {
        if(token.trim() != "") {
            pushCatagorizedToken(token.trim());
        }

        token = "";

        if(also != undefined) {
            if(also.trim() != "") {
                pushCatagorizedToken(also);
            }
        }

        tokenStart = i;
    }

    for(var i = 0; i < source.length; i++) {
        let char = source[i];

        if(char == "\n") { line++; }

        if(char == '"') {
            in_quote = !in_quote;

            if(in_quote) {
                pushToken();
                token += char;
                continue
            } else {
                token += char;
                pushToken();
                continue
            }
        }

        if(in_quote) {
            token += char;
            continue
        }
 
        char_type = getCharType(char);

        // console.log(char_type);

        if(i == 0) { last_char_type = char_type; }

        if(char_type != last_char_type || char_type == "other") {
            pushToken(char_type == "other" ? char : undefined);
            if(char_type != "other") { token += char; }
        } else {
            token += char;
        }

        last_char_type = char_type;
    }

    pushToken();

    return tokens;
}