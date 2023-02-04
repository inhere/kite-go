/* https://github.com/Khan/KaTeX/issues/350#issuecomment-160906578 */
.markdown-body {
    counter-reset: eqnum;
}

.katex-display {
    display: table;
    width: 100%;
}

.katex-display>.katex {
    display: table-cell;
}

.katex-display::before,
.katex-display::after {
    width: 10000px;
    display: table-cell;
    text-align: right;
    vertical-align: middle;
}

.katex-display::after {
    counter-increment: eqnum;
    content: "(" counter(eqnum) ")";
}
