package alarm

type Type string

const (
	None       Type = "none"
	EntryDelay Type = "entry-delay"
	Burglar    Type = "burglar"
	Panic      Type = "panic"
	Fire       Type = "fire"
	Tamper     Type = "tamper"
	TamperPin  Type = "tamper-pin"
)

func (a Type) IsValid() bool {
	return a == None || a == Burglar || a == Panic || a == Fire || a == Tamper || a == TamperPin || a == EntryDelay
}
