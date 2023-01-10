"use strict"

const _ = undefined,

    element_constructor = (element, className, idName, innerTxt, srcName, innerHTML) => {
        let el = document.createElement(element);

        className != _ ? el.setAttribute("class", className) : null;
        idName != _ ? el.setAttribute("id", idName) : null;
        innerTxt != _ ? el.innerText = innerTxt : null;
        srcName != _ ? (_ => { el.setAttribute("src", srcName); el.setAttribute("alt", srcName) })() : null;
        innerHTML != _ ? el.innerHTML = innerHTML : null;

        return el;
    };

export { element_constructor as default, _ };