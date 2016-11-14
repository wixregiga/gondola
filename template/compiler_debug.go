package template

import (
	"fmt"
)

func (o opcode) String() string {
	switch o {
	case opDOT:
		return "opDOT"
	case opFIELD:
		return "opFIELD"
	case opFUNC:
		return "opFUNC"
	case opITER:
		return "opITER"
	case opJMP:
		return "opJMP"
	case opJMPF:
		return "opJMPF"
	case opJMPT:
		return "opJMPT"
	case opMARK:
		return "opMARK"
	case opNEXT:
		return "opNEXT"
	case opCONTEXT:
		return "opCONTEXT"
	case opNOP:
		return "opNOP"
	case opPOP:
		return "opPOP"
	case opPOPDOT:
		return "opPOPDOT"
	case opPRINT:
		return "opPRINT"
	case opPUSHDOT:
		return "opPUSHDOT"
	case opSETVAR:
		return "opSETVAR"
	case opSTATE:
		return "opSTATE"
	case opSTRING:
		return "opSTRING"
	case opTEMPLATE:
		return "opTEMPLATE"
	case opUNSETVAR:
		return "opUNSETVAR"
	case opVAL:
		return "opVAL"
	case opVAR:
		return "opVAR"
	case opWB:
		return "opWB"
	}
	return fmt.Sprintf("unknown opcode %d", o)
}

func (p *program) dumpTemplate(tmpl string) {
	ins := p.code[tmpl]
	fmt.Printf("Template %q: %d instructions\n", tmpl, len(ins))
	dumpInstructions(p, ins)
}

func dumpInstructions(p *program, ins []inst) {
	for ii, v := range ins {
		var value string
		switch v.op {
		case opFIELD:
			args, i := decodeVal(v.val)
			value = fmt.Sprintf("%q - %d args", p.strings[i], args)
		case opFUNC:
			args, i := decodeVal(v.val)
			value = fmt.Sprintf("%q - %d args", p.funcs[i].Name, args)
		case opSETVAR, opUNSETVAR, opVAR, opSTRING:
			value = fmt.Sprintf("%q", p.strings[int(v.val)])
		case opTEMPLATE:
			n, t := decodeVal(v.val)
			name := p.strings[t]
			ns := p.strings[n]
			value = fmt.Sprintf("%q - %q", name, ns)
		case opVAL:
			value = fmt.Sprintf("%+v", p.values[v.val].Interface())
		case opWB:
			b := p.bs[int(v.val)]
			value = fmt.Sprintf("%d - %d bytes", v.val, len(b))
		default:
			im := int(int32(v.val))
			if im != 0 || v.op == opPOP {
				value = fmt.Sprintf("%d", im)
			}
		}
		fmt.Printf("PC %d %s %v\n", ii, v.op, value)
	}
}
