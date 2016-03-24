package lazyblog

var (
	css = []byte(`
html, body {
  font-family: 'Inconsolata', Menlo, monospace;
  margin: 0;
  color: rgba(0,0,0,.8);
  font-size: 1.0625rem;
  line-height: 1.75;
}

h1 {
  font-size: 1.5rem;
  text-transform: uppercase;
  margin-top: 2rem;
  margin-bottom: 1.75rem;
  text-align: center;
}

h3 {
  font-size: 1.25rem;
  text-transform: uppercase;
}

p {
  font-size: 1rem;
}

input,
textarea {
  border: none;
  outline: none;
  font-size: 1rem;
  font-family: 'Inconsolata', Menlo, monospace;
  width: 100%;
}

textarea {
  height: auto;
  resize: none;
}

.writing {
  margin-top: 2rem;
  margin-bottom: 2rem;
}

.wrapper {
  max-width: 728px;
  margin: 0 auto;
}

.btn {
  border: none;
  padding: 6px 12px;
  font-size: .75rem;
  text-transform: uppercase;
  background: rgba(0,0,0,.8);
  color: #fff;
  width: auto;
  margin-bottom: 2.5rem;
}

.btn > a {
  color: #fff; text-decoration: none;
}

.btn:hover {
  cursor: pointer;
}

.btn:visited {
  color: #fff;
}

.px2 {
  padding-left: 1rem; padding-right: 1rem;
}
		`)
)
