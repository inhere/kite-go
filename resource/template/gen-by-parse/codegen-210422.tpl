# usage: kit cgen tmp/codegen-210422.tpl
Levels = [Info,Trace,Error,Notice,Debug,Fatal,Panic]
Name = r
Struct = *Record

###

{{range $lv := .Levels}}
// {{$lv}} logs a message at level {{ $lv }}
func ({{$.Name}} {{$.Struct}}) {{ $lv }}(args ...interface{}) {
	{{ $lv }}.Log({{ $lv }}Level, args...)
}
{{end}}
