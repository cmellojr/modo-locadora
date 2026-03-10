package almanac

import "time"

// ephemerides maps day-of-year (1-indexed) to a gaming curiosity in Portuguese.
var ephemerides = map[int]string{
	// Janeiro
	1:  "Em 1985, o Famicom chegava ao ocidente como Nintendo Entertainment System.",
	15: "Em 1990, o Neo Geo AES era lancado como o console mais caro da historia.",
	25: "Em 1995, o Sega Saturn era lancado no Japao.",
	31: "Em 1998, Resident Evil 2 era lancado no PlayStation e se tornava um classico do terror.",

	// Fevereiro
	35: "Em 1986, The Legend of Zelda era lancado no Famicom Disk System no Japao.",
	42: "Em 1989, o Game Boy era anunciado pela Nintendo.",
	50: "Em 2005, o Nintendo DS ja vendia mais de 5 milhoes de unidades.",

	// Marco
	60: "Em 1983, o Famicom era lancado no Japao, iniciando a era moderna dos consoles.",
	68: "Em 1994, Super Metroid chegava ao SNES e redefiniu o genero de exploracao.",
	75: "Em 1987, Mega Man era lancado no Famicom, iniciando uma das maiores sagas.",
	80: "Em 1991, Streets of Rage chegava ao Mega Drive.",

	// Abril
	92:  "Em 1989, o Game Boy era lancado no Japao com Tetris incluido.",
	100: "Em 1992, Sonic 2 estava em desenvolvimento nos EUA pelo Sega Technical Institute.",
	105: "Em 1999, o Dreamcast recebia Soul Calibur, um dos melhores jogos de luta.",

	// Maio
	121: "Em 1989, o Mega Drive chegava ao Brasil como primeiro console de 16 bits.",
	130: "Em 1993, Star Fox usava o chip Super FX para graficos 3D no SNES.",
	140: "Em 2001, o GameCube era apresentado na E3.",

	// Junho
	152: "Em 1991, Sonic the Hedgehog estreava no Mega Drive e desafiava o Mario.",
	160: "Em 1996, o Nintendo 64 era apresentado na E3 com Super Mario 64.",
	168: "Em 1995, o Sega Saturn era lancado de surpresa na E3 americana.",
	175: "Em 2004, o Nintendo DS era revelado na E3.",

	// Julho
	183: "Em 1997, Final Fantasy VII era lancado no Japao e mudava os RPGs para sempre.",
	190: "Em 1993, Doom era lancado e revolucionava os jogos de PC.",
	200: "Em 1990, o Super Famicom era anunciado no Japao.",

	// Agosto
	213: "Em 1988, Mega Man 2 era lancado no Famicom — o mais popular da serie.",
	220: "Em 1992, Mortal Kombat causava polemica mundial nos fliperamas.",
	230: "Em 1987, Castlevania estreava no NES na America do Norte.",

	// Setembro
	244: "Em 1998, The Legend of Zelda: Ocarina of Time se preparava para o lancamento.",
	250: "Em 1999, o Dreamcast era lancado nos EUA, primeiro console com modem.",
	258: "Em 1994, o PlayStation era lancado no Japao e mudava tudo.",

	// Outubro
	274: "Em 1985, Super Mario Bros. era lancado no Japao e salvava a industria.",
	282: "Em 1988, Phantasy Star abria o caminho dos RPGs no Master System.",
	290: "Em 1986, Dragon Quest inaugurava o genero RPG nos consoles.",

	// Novembro
	305: "Em 1990, o Super Famicom era lancado no Japao.",
	310: "Em 1994, Donkey Kong Country impressionava com graficos pre-renderizados no SNES.",
	318: "Em 1998, The Legend of Zelda: Ocarina of Time recebia notas perfeitas.",
	324: "Em 1992, Sonic 2 era lancado no 'Sonic 2sday' em todo o mundo.",

	// Dezembro
	336: "Em 1996, o Nintendo 64 chegava as lojas com Super Mario 64.",
	345: "Em 1987, Contra no NES popularizou o codigo Konami.",
	355: "Em 1995, Chrono Trigger era lancado no Super Famicom.",
	360: "Em 1994, Super Metroid ja era considerado um dos melhores jogos de todos os tempos.",
}

// genericEphemerides are used for days without a specific entry.
var genericEphemerides = []string{
	"Voce sabia? O NES vendeu mais de 61 milhoes de unidades no mundo.",
	"Voce sabia? O Mega Drive foi o console mais popular do Brasil nos anos 90.",
	"Voce sabia? O codigo Konami apareceu em mais de 100 jogos diferentes.",
	"Voce sabia? Pac-Man foi o primeiro personagem de videogame a virar febre mundial.",
	"Voce sabia? O Game Boy sobreviveu a uma explosao na Guerra do Golfo e continuou funcionando.",
	"Voce sabia? Mario foi batizado em homenagem ao senhorio da Nintendo of America.",
	"Voce sabia? Tetris foi criado por um programador sovietico em 1984.",
	"Voce sabia? O SNES tinha um chip de som projetado pela Sony.",
	"Voce sabia? O Atari 2600 vendeu mais de 30 milhoes de unidades.",
	"Voce sabia? Shigeru Miyamoto quase chamou Zelda de 'Hyrule Fantasy'.",
	"Voce sabia? O Game Boy Color era retrocompativel com todos os jogos de Game Boy.",
	"Voce sabia? Street Fighter II e o jogo de luta mais influente da historia.",
	"Voce sabia? A TecToy adaptou mais de 100 jogos exclusivos para o Master System no Brasil.",
	"Voce sabia? O Mega Drive teve suporte oficial no Brasil ate 2002.",
	"Voce sabia? A Playtronic foi a distribuidora oficial da Nintendo no Brasil nos anos 90.",
	"Voce sabia? O primeiro console a usar CD-ROM foi o TurboGrafx-CD em 1988.",
	"Voce sabia? Lara Croft de Tomb Raider foi um dos primeiros icones femininos dos games.",
	"Voce sabia? O Neo Geo custava o equivalente a um computador quando foi lancado.",
	"Voce sabia? A franquia Pokemon ja vendeu mais de 440 milhoes de jogos.",
	"Voce sabia? O primeiro easter egg em videogames foi em Adventure, no Atari 2600.",
}

// TodaysEphemeride returns the gaming curiosity for today.
func TodaysEphemeride() string {
	dayOfYear := time.Now().YearDay()
	if text, ok := ephemerides[dayOfYear]; ok {
		return text
	}
	return genericEphemerides[dayOfYear%len(genericEphemerides)]
}
