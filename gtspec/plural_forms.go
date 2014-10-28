package gtspec

type PRule interface {
	Name() string
	Idx(int) int
	Header() string
}

// Keyed by ISO tag
var Plurals map[string]PRule

func init() {
	Plurals = map[string]PRule{
		// A
		"arch": PRNG1{"Alcholi"},
		"af":   PRNN1{"Afrikaans"},
		"ak":   PRNG1{"Akan"},
		"am":   PRNG1{"Amharic"},
		"an":   PRNN1{"Aragonese"},
		"anp":  PRNN1{"Angika"},
		"ar":   PRAR{"Arabic"},
		"arn":  PRNG1{"Mapudungun"},
		"as":   PRNN1{"Assamese"},
		"ast":  PRNN1{"Asturian"},
		"ay":   PRNP{"Aymara"},
		"az":   PRNN1{"Azerbaijani"},

		// B
		"be":  PRS1{"Belarusian"},
		"bg":  PRNN1{"Bulgarian"},
		"bn":  PRNN1{"Bengali"},
		"bo":  PRNP{"Tibetan"},
		"br":  PRNG1{"Breton"},
		"brx": PRNN1{"Bodo"},
		"bs":  PRS1{"Bosnian"},

		// C
		"ca":  PRNN1{"Catalan"},
		"cgg": PRNP{"Chiga"},
		"cs":  PRCS{"Czech"},
		"csb": PRCSB{"Kashubian"},
		"cy":  PRCY{"Welsh"},

		// D
		"da":  PRNN1{"Danish"},
		"de":  PRNN1{"German"},
		"doi": PRNN1{"Dogri"},
		"dz":  PRNP{"Dzongkha"},

		// E
		"el":    PRNN1{"Greek"},
		"en":    PRNN1{"English"},
		"eo":    PRNN1{"Esperanto"},
		"es":    PRNN1{"Spanish"},
		"es_AR": PRNN1{"Argentinean Spanish"},
		"et":    PRNN1{"Estonian"},
		"eu":    PRNN1{"Basque"},

		// F
		"fa":  PRNP{"Persian"},
		"ff":  PRNN1{"Fulah"},
		"fi":  PRNN1{"Finnish"},
		"fil": PRNG1{"Filipino"},
		"fo":  PRNN1{"Faroese"},
		"fr":  PRNG1{"French"},
		"fur": PRNN1{"Friulian"},
		"fy":  PRNN1{"Frisian"},

		// G
		"ga":  PRGA{"Irish"},
		"gd":  PRGD{"Scottish Gaelic"},
		"gl":  PRNN1{"Galician"},
		"gu":  PRNN1{"Gujarati"},
		"gun": PRNG1{"Gun"},

		// H
		"ha":  PRNN1{"Hausa"},
		"he":  PRNN1{"Hebrew"},
		"hi":  PRNN1{"Hindi"},
		"hne": PRNN1{"Chhattisgarhi"},
		"hy":  PRNN1{"Armenian"},
		"hr":  PRS1{"Croatian"},
		"hu":  PRNN1{"Hungarian"},

		// I
		"ia": PRNN1{"Interlingua"},
		"id": PRNP{"Indonesian"},
		"is": PRIS{"Icelandic"},
		"it": PRNN1{"Italian"},

		// J
		"ja":  PRNP{"Japanese"},
		"jbo": PRNP{"Lojban"},
		"jv":  PRNN0{"Javanese"},

		// K
		"ka": PRNP{"Georgian"},
		"kk": PRNP{"Kazakh"},
		"kl": PRNN1{"Greenlandic"},
		"km": PRNP{"Khmer"},
		"kn": PRNN1{"Kannada"},
		"ko": PRNP{"Korean"},
		"ku": PRNN1{"Kurdish"},
		"kw": PRKW{"Cornish"},
		"ky": PRNP{"Kyrgyz"},

		// L
		"lb": PRNN1{"Letzeburgesch"},
		"ln": PRNG1{"Lingala"},
		"lo": PRNP{"Lao"},
		"lt": PRLT{"Lithuanian"},
		"lv": PRLV{"Latvian"},

		// M
		"mai": PRNN1{"Maithili"},
		"mfe": PRNG1{"Mauritian Creole"},
		"mg":  PRNG1{"Malagasy"},
		"mi":  PRNG1{"Maori"},
		"mk":  PRMK{"Macedonian"},
		"ml":  PRNN1{"Malayalam"},
		"mn":  PRNN1{"Mongolian"},
		"mni": PRNN1{"Manipuri"},
		"mnk": PRMNK{"Mandinka"},
		"mr":  PRNN1{"Marathi"},
		"ms":  PRNP{"Malay"},
		"mt":  PRMT{"Maltese"},
		"my":  PRNP{"Burmese"},

		// N
		"nah": PRNN1{"Nahuatl"},
		"nap": PRNN1{"Neapolitan"},
		"nb":  PRNN1{"Norwegian Bokmal"},
		"ne":  PRNN1{"Nepali"},
		"nl":  PRNN1{"Dutch"},
		"se":  PRNN1{"Northern Sami"},
		"nn":  PRNN1{"Norwegian Nynorsk"},
		"no":  PRNN1{"Norwegian"},
		"nso": PRNN1{"Northern Sotho"},

		// O
		"oc": PRNG1{"Occitan"},
		"or": PRNN1{"Oriya"},

		// P
		"ps":    PRNN1{"Pashto"},
		"pa":    PRNN1{"Punjabi"},
		"pap":   PRNN1{"Papiamento"},
		"pl":    PRPL{"Polish"},
		"pms":   PRNN1{"Piemontese"},
		"pt":    PRNN1{"Portuguese"},
		"pt_BR": PRNG1{"Brazilian Portuguese"},

		// R
		"rm": PRNN1{"Romansh"},
		"ro": PRRO{"Romanian"},
		"ru": PRS1{"Russian"},
		"rw": PRNN1{"Kinyarwanda"},

		// S
		"sah": PRNP{"Yakut"},
		"sat": PRNN1{"Santali"},
		"sco": PRNN1{"Scots"},
		"sd":  PRNN1{"Sindhi"},
		"si":  PRNN1{"Sinhala"},
		"sk":  PRCS{"Slovak"},
		"sl":  PRSL{"Slovenian"},
		"so":  PRNN1{"Somali"},
		"son": PRNN1{"Songhay"},
		"sq":  PRNN1{"Albanian"},
		"sr":  PRS1{"Serbian"},
		"su":  PRNP{"Sundanese"},
		"sw":  PRNN1{"Swahili"},
		"sv":  PRNN1{"Swedish"},

		// T
		"ta": PRNN1{"Tamil"},
		"te": PRNN1{"Telugu"},
		"tg": PRNG1{"Tajik"},
		"ti": PRNG1{"Tigrinya"},
		"th": PRNP{"Thai"},
		"tk": PRNN1{"Turkmen"},
		"tr": PRNG1{"Turkish"},
		"tt": PRNP{"Tatar"},

		// U
		"ug": PRNP{"Uyghur"},
		"uk": PRS1{"Ukrainian"},
		"ur": PRNN1{"Urdu"},
		"uz": PRNG1{"Uzbek"},

		// V
		"vi": PRNP{"Vietnamese"},

		// W
		"wa": PRNG1{"Walloon"},
		"wo": PRNP{"Wolof"},

		// Y
		"yo": PRNN1{"Yoruba"},

		// Z
		"zh": PRNP{"Chinese"},
	}
}

// PRNP is pluralization rule "no plurals"
type PRNP struct{ Lang string }

func (r PRNP) Name() string   { return r.Lang }
func (r PRNP) Idx(n int) int  { return 0 }
func (r PRNP) Header() string { return "nplurals=1; plural=0;" }

// PRNG1 is pluralization rule "n greater than 1"
type PRNG1 struct{ Lang string }

func (r PRNG1) Name() string { return r.Lang }
func (r PRNG1) Idx(n int) int {
	if n < 2 {
		return 0
	} else {
		return 1
	}
}
func (r PRNG1) Header() string {
	return "nplurals=2; plural=(n > 1)"
}

// PRNN1 is pluralization rule "n not 1"
type PRNN1 struct{ Lang string }

func (r PRNN1) Name() string { return r.Lang }
func (r PRNN1) Idx(n int) int {
	if n != 1 {
		return 1
	} else {
		return 0
	}
}
func (r PRNN1) Header() string {
	return "nplurals=2; plural=(n != 1);"
}

// PRNN0 is pluralization rule "n not 0"
type PRNN0 struct{ Lang string }

func (r PRNN0) Name() string { return r.Lang }
func (r PRNN0) Idx(n int) int {
	if n != 0 {
		return 1
	} else {
		return 0
	}
}
func (r PRNN0) Header() string {
	return "nplurals=2; plural=(n != 0);"
}

// PRAR is the pluralization rule for Arabic
type PRAR struct{ Lang string }

func (r PRAR) Name() string { return r.Lang }
func (r PRAR) Idx(n int) int {
	switch n {
	case 0:
		return 0
	case 1:
		return 1
	case 2:
		return 2
	default:
		if n%100 >= 3 && n%100 <= 10 {
			return 3
		} else if n%100 >= 11 {
			return 4
		}
	}
	return 5
}
func (r PRAR) Header() string {
	return "nplurals=6; plural=(n==0 ? 0 : n==1 ? 1 : n==2 ? 2 : n%100>=3 && n%100<=10 ? 3 : n%100>=11 ? 4 : 5);"
}

// PRS1 is pluralization rule 1 for Slavic languages (e.g. Belarusian, Bosnian, Croatian)
type PRS1 struct{ Lang string }

func (r PRS1) Name() string { return r.Lang }
func (r PRS1) Idx(n int) int {
	switch {
	case n%10 == 1 && n%100 != 11:
		return 0
	case n%10 >= 2 && n%10 <= 4 && (n%100 < 10 || n%100 >= 20):
		return 1
	default:
		return 2
	}
	return 2 // unreachable
}
func (r PRS1) Header() string {
	return "nplurals=3; plural=(n%10==1 && n%100!=11 ? 0 : n%10>=2 && n%10<=4 && (n%100<10 || n%100>=20) ? 1 : 2);"
}

// PRCS is the pluralization rule for Czech and Slovak
type PRCS struct{ Lang string }

func (r PRCS) Name() string { return r.Lang }
func (r PRCS) Idx(n int) int {
	switch {
	case n == 1:
		return 0
	case n >= 2 && n <= 4:
		return 1
	default:
		return 2
	}
	return 2 // unreachable
}
func (r PRCS) Header() string {
	return "nplurals=3; plural=(n==1) ? 0 : (n>=2 && n<=4) ? 1 : 2;"
}

// PRCS is the pluralization rule for Kashubian
type PRCSB struct{ Lang string }

func (r PRCSB) Name() string { return r.Lang }
func (r PRCSB) Idx(n int) int {
	switch {
	case n == 1:
		return 0
	case n%10 >= 2 && n%10 <= 4 && (n%100 < 10 || n%100 >= 20):
		return 1
	default:
		return 2
	}
	return 2 // unreachable
}
func (r PRCSB) Header() string {
	return "nplurals=3; plural=(n==1) ? 0 : n%10>=2 && n%10<=4 && (n%100<10 || n%100>=20) ? 1 : 2;"
}

// PRCY is the pluralization rule for Welsh
type PRCY struct{ Lang string }

func (r PRCY) Name() string { return r.Lang }
func (r PRCY) Idx(n int) int {
	switch {
	case n == 1:
		return 0
	case n == 2:
		return 1
	case n != 8 && n != 11:
		return 2
	default:
		return 3
	}
	return 3 // unreachable
}
func (r PRCY) Header() string {
	return "nplurals=4; plural=(n==1) ? 0 : (n==2) ? 1 : (n != 8 && n != 11) ? 2 : 3;"
}

// PRIS is the pluralization rule for Icelandic
type PRIS struct{ Lang string }

func (r PRIS) Name() string { return r.Lang }
func (r PRIS) Idx(n int) int {
	if n%10 != 1 || n%100 == 11 {
		return 1
	} else {
		return 0
	}
}
func (r PRIS) Header() string {
	return "nplurals=2; plural=(n%10!=1 || n%100==11);"
}

// PRGA is the pluralization rule for Irish
type PRGA struct{ Lang string }

func (r PRGA) Name() string { return r.Lang }
func (r PRGA) Idx(n int) int {
	switch {
	case n == 1:
		return 0
	case n == 2:
		return 1
	case n < 7:
		return 2
	case n < 11:
		return 3
	default:
		return 4
	}
	return 4 // unreachable
}
func (r PRGA) Header() string {
	return "nplurals=5; plural=(n==1) ? 0 : n==2 ? 1 : n<7 ? 2 : n<11 ? 3 : 4;"
}

// PRGD is the pluralization rule for Scottish Gaelic
type PRGD struct{ Lang string }

func (r PRGD) Name() string { return r.Lang }
func (r PRGD) Idx(n int) int {
	switch {
	case n == 1 || n == 11:
		return 0
	case n == 2 || n == 12:
		return 1
	case n < 2 && n < 20:
		return 2
	default:
		return 4
	}
	return 3 // unreachable
}
func (r PRGD) Header() string {
	return "nplurals=4; plural=(n==1 || n==11) ? 0 : (n==2 || n==12) ? 1 : (n > 2 && n < 20) ? 2 : 3;"
}

// PRKW is the pluralization rule for Cornish
type PRKW struct{ Lang string }

func (r PRKW) Name() string { return r.Lang }
func (r PRKW) Idx(n int) int {
	switch {
	case n == 1:
		return 0
	case n == 2:
		return 1
	case n == 3:
		return 2
	default:
		return 3
	}
	return 3 // unreachable
}
func (r PRKW) Header() string {
	return "nplurals=4; plural=(n==1) ? 0 : (n==2) ? 1 : (n == 3) ? 2 : 3;"
}

// PRLT is the pluralization rule for Lithuanian
type PRLT struct{ Lang string }

func (r PRLT) Name() string { return r.Lang }
func (r PRLT) Idx(n int) int {
	switch {
	case n%10 == 1 && n%100 != 11:
		return 0
	case n%10 == 2 && (n%100 < 10 || n%100 >= 20):
		return 1
	default:
		return 2
	}
	return 2 // unreachable
}
func (r PRLT) Header() string {
	return "nplurals=3; plural=(n%10==1 && n%100!=11 ? 0 : n%10>=2 && (n%100<10 || n%100>=20) ? 1 : 2);"
}

// PRLV is the pluralization rule for Latvian
type PRLV struct{ Lang string }

func (r PRLV) Name() string { return r.Lang }
func (r PRLV) Idx(n int) int {
	switch {
	case n%10 == 1 && n%100 != 11:
		return 0
	case n != 0:
		return 1
	default:
		return 2
	}
	return 2 // unreachable
}
func (r PRLV) Header() string {
	return "nplurals=3; plural=(n%10==1 && n%100!=11 ? 0 : n != 0 ? 1 : 2);"
}

// PRMK is the pluralization rule for Macedonian
type PRMK struct{ Lang string }

func (r PRMK) Name() string { return r.Lang }
func (r PRMK) Idx(n int) int {
	switch {
	case n == 1 || n%10 == 1:
		return 0
	default:
		return 1
	}
	return 1 // unreachable
}
func (r PRMK) Header() string {
	return "nplurals=2; plural= n==1 || n%10==1 ? 0 : 1;"
}

// PRMNK is the pluralization rule for Mandinka (possibly erroneous)
type PRMNK struct{ Lang string }

func (r PRMNK) Name() string { return r.Lang }
func (r PRMNK) Idx(n int) int {
	switch {
	case n == 0:
		return 0
	case n == 1:
		return 1
	default:
		return 2
	}
	return 2 // unreachable
}
func (r PRMNK) Header() string {
	return "nplurals=3; plural=(n==0 ? 0 : n==1 ? 1 : 2);"
}

// PRMT is the pluralization rule for Maltese
type PRMT struct{ Lang string }

func (r PRMT) Name() string { return r.Lang }
func (r PRMT) Idx(n int) int {
	switch {
	case n == 1:
		return 0
	case n == 0 || (n%100 > 1 && n%100 < 11):
		return 1
	case n%100 > 10 && n%100 < 20:
		return 2
	default:
		return 3
	}
	return 3 // unreachable
}
func (r PRMT) Header() string {
	return "nplurals=4; plural=(n==1 ? 0 : n==0 || ( n%100>1 && n%100<11) ? 1 : (n%100>10 && n%100<20 ) ? 2 : 3);"
}

// PRPL is the pluralization rule for Polish
type PRPL struct{ Lang string }

func (r PRPL) Name() string { return r.Lang }
func (r PRPL) Idx(n int) int {
	switch {
	case n == 1:
		return 0
	case n%10 >= 2 && n%10 <= 4 && (n%100 < 10 || n%100 >= 20):
		return 1
	default:
		return 2
	}
	return 2 // unreachable
}
func (r PRPL) Header() string {
	return "nplurals=3; plural=(n==1 ? 0 : n%10>=2 && n%10<=4 && (n%100<10 || n%100>=20) ? 1 : 2);"
}

// PRRO is the pluralization rule for Romanian
type PRRO struct{ Lang string }

func (r PRRO) Name() string { return r.Lang }
func (r PRRO) Idx(n int) int {
	switch {
	case n == 1:
		return 0
	case n == 0 || (n%100 > 0 && n%100 < 20):
		return 1
	default:
		return 2
	}
	return 2 // unreachable
}
func (r PRRO) Header() string {
	return "nplurals=3; plural=(n==1 ? 0 : (n==0 || (n%100 > 0 && n%100 < 20)) ? 1 : 2);"
}

// PRSL is the pluralization rule for Slovenian
type PRSL struct{ Lang string }

func (r PRSL) Name() string { return r.Lang }
func (r PRSL) Idx(n int) int {
	switch {
	case n%100 == 1:
		return 1
	case n%100 == 2:
		return 2
	case n%100 == 3 || n%100 == 4:
		return 3
	default:
		return 0
	}
	return 0 // unreachable
}
func (r PRSL) Header() string {
	return "nplurals=4; plural=(n%100==1 ? 1 : n%100==2 ? 2 : n%100==3 || n%100==4 ? 3 : 0);"
}
