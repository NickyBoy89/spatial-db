LATEX := latexmk -pdf --shell-escape

all: document.pdf

document.pdf: document.tex references.bib oxycomps.sty
	$(LATEX) document.tex

.PHONY: clean
clean:
	$(LATEX) -C
	rm *.bbl
