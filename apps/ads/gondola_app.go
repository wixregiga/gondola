package ads

// AUTOMATICALLY GENERATED WITH gondola gen-app -release -- DO NOT EDIT!

import (
	"gnd.la/app"
	"gnd.la/loaders"
	"gnd.la/template"
	"gnd.la/template/assets"
)

var _ = loaders.FSLoader
var _ = template.New
var _ = assets.NewManager
var (
	App = app.New()
)

func init() {
	App.SetName("Ads")
	var manager *assets.Manager
	assetsLoader := loaders.TgzLoader(loaders.RawString("\x1f\x8b\b\x00\x00\tn\x88\x02\xff\xec\x1a]o\xe36r\x9f\xf3+\xb8>_%#\xb6\xd6\xcenv{\xf6\xa5hn\x81\xa2\xc1\x01y\xb8ݷ X\xd0\x12e\xb3\x91EU\xa4\x9c\xa4\xa9\xff\xfb͐Էd;Ew\xef\x0e\x17bQ\x8b\xe4|\xcfpfȔ\x06\xd2\xf3\xa5|\xf55\xc7Ԍ\xbe߳\xf3\x0fo\xf3o\xb3>\x9b\xbe\xfd\xf0\xe1\x15y\xf5-F&\x15M\x81\xfd\x9f\xad\xe4\xff\xc8\xf0 \x00&J$\x93\xa5x O'\x04Ɔ\xa6+\x1e\xcfɔ\xd0L\x89\x85^K\x84\xe4\x8a\vXMYD\x15߲\xc5\xc9\xeeD#\x97\x88K\xea߭R\x91\xc5\xc1\x9c\xfc%\fC\x83\xaaR\x1a\xe7\xc8\x00;\x91k\x1a\x88{2\xf5\xdeJ¨d\x13\x1eOD\xa6\x90^Ap\xb2\xe5\x92/#F\xf4\x82ŰL\n\x12(\xe1\x94\xccΓ\a\x92\xae\x96ԝ\x8e\x89\xfd睝\x8fJz)\x93\x89\x00\x11\xb6\x96\\)\xf0=\x0f\xd4zNޞM\x93\a#\xec\x9a\xf1\xd5Z\xcdɹ^\x01\n?nX\xc0\xa9\xbb\x01!-\xf4\xf9\x14\xf6F\x96\xc2\x01\x06\x15&\xef\xde\u007f\x9f3\xa92z_\xb0\xde\x01\xbb6\xb7\xef\xff\x18\xb7\x0fg\x9d\xdc\xfeV\xe3V\xfaO\u007f\xacy\xc0&\xcbL)\x11[\x82\x01\x97ID\x1f\xe7$\x161k\x06\x02]J\x11eʮG,\x04\xfaS3\xf9\r|\x1a\xb0\a\xb0\xac=\r\x16\x97\x06\x01\x8fW\xe8\xb7wM\x83\xcf\xde\xe7+\xb5(*p}\x11\x89\xb4\x1aV\xa1\x88\xd5D\xf2\xdf\x18\xe0\x9eYg5\xed\x83Z\x14L1N41k\xa2\xd9t\xfaW\x9c\x16\xc2ά\xa89\xa1\x90?\xb0\xa0NC\xe3\x97&\xd0\x10-\x925|<[\x9a\x06\xfc\xce\xc9\x04yT\x85-az\\\xb0\x140\xd9 fi \x91\x06,\x9d\x98\x9d\t\x1a~\x92Ҁg\x12\xec\xdd\r\x93\xa2\x8d\xeb@u\x01\f\x9cfY2\xec\x14\xd5B\xf6Hk\x95l\x8a\x8a\xf9e\xaf\x9c\b\xd0-d\xd5\x15\xad\\\xd5\xccT\xb9\xa8(ؘ\xf4%\x13\x9bJ\x1a\xc1\x9d\xe3\xe2~\x1f.ҭ\xe3\xf28\xe2q\x99\vY\x9a\x8aTC\xd8xMM|\xe8`\xbd\xb7\xa1\xbe\x14Q\x80\x18\xff\x15\xf9\x1f\xeb\xff/_\xb7\xfc\x1f\xa8\xff\xb3w\xeffg\xcd\xfa\u007f6{\xffR\xff\xbf\xc5p\xc3,\xf61\x9b\xb91ݰ\xbc\xd0liJbI. \xb1\xc5Piop\xef\xb69\xfd\xfdw\xf2\xb43\xc7\xf8\xe9\x89\f#*\x15\x99_\x10ȕčXL~LR\xb1\x85\x13\x93\xca\x11\xe4\x04\xb2\xdb\x15\xa4\x03\x16\xd2,RȠ\xac\\@\x03:\x85\x15#C\xce\xc7d\xb8EZ%\x89\x1c\xbd\x02=\xdcz\xd7 \b\xec̡\xc9\xc8\x18,\xf1\x90\xc4\xecW\xa4`\xe5\xd9\xedư\xcc⠊__\x01\x15J\x9d\x99\xba\x17\xe9ݟ#\x18\xca\xf5T\xdbÁ;s\xe2\xd4\xe4w\xc6-\xb0\x8f \xbe,\xe0|\x9c\xf5\x03\u007f\xf2S\x9e\xa8\x02Z\xea)\x82\xd6 w\x87\fD\x0eX(\x96\xdeU\xcc\x15\x18\xa7\b\x1a\x91\xe0\x8f\x1cU\x14\xb5K\x005\xf4\u0603\x02:\xee\xd3n\\\xf8|\x9c\x03\x8c\xca\xf6\x04\xa4r\x81\xf8\x97\x94\xd1\xe0q\xd40Z\xcaT\x96\xc6%pi\xee\x02\x05xa\x00\x940oސ\x8fk\xe6\xdf!\xe5{F\xd6\x14\x9a\x01\x1a?BlJ\x02\xc5J\xad\x19\xd4\xf4\x15\xab\t0t\x9d<\xeb;#\x0f\x02x\xa5\xd6\xe4\xe2\x82L\x8f\x17(\x84\xe4\xefb\x1cݱG(\rE85)TB\rDϡn\x00\xebv\xd1\x02\x94,b\xbe\x02\xca\x17\xc4\xf1\x1cr\x9a\xc3{:D\xda\xf0\xa8#\xd8\xde\xcd\xf1Fu\x10T\x15\x93\xbeU\xf0\x87\xb6~8\x10\x82Q\u007f]\xa4\a\xe2v\x81\xe5<\x87j\xcd\rW\xfchp\xacr~\xad!=.]gnK\xab3\ua8dcS\x0f\xd1H\xe0\xeb\x15S\xd7F\xf7\x9f\x8a\xa4e\xe6c\xe2\xa4l#\xb6@l\xd1K\n\xf9\x87\xfbxi\x17\xba\xc3=*\xd4\xdd\xdd\x1cF7#\x88\xdbC\xa1\x8d\xbd\xeb\x80<\xe0\xc3\xdd\x1f\xf0\xa8\t%\x95%\xc7\x18S\x03v\xd9\x12y\xe9\xcd>;>#pʀ=\x14984WгϮ\x1d\xeb\xbbN\xf1m\xfaѧ\xadO4\x9d\x90\xec13\xc9u\x9f\x12e\r-\xfcu \xca|\x10@D̋\xc4\xcau\"A\xf1vA\x9c1i\xf0\\\xec%2\xf4\xc0\x87\x06\xb2!\xec\xb8L\xd0\x01UtL\x14$\xe2O\x8a\xaa\f\xd2\xef/\xbf>\xac\x0f\n\x98\x1bA\x03{R\xa3b2<\x9bN\x8fAͳ\xf3'h\xa6/\x03\xd9\x17ǝ\x8e$,\x92\xec\x19<~\x862\xfc\\\x1e{!v\xfb\x0e\xff\xe8\x10\x9f\x83\xf2?\xc3.\x1d\xc9b\x1f\xf5#\xad\xb1\xeb\xc9%\xb6\xd8\x17\xb5>\x97\xb2Z\xef\xbb\x02</\x9c\xfa\x92\x06\xa5\xb3v\xfcݮڷ\xb7\\\xe82Q\x14=h2MZ\r\xa1\xf5,\xd9\xdbd\xd7\x15\x8b\xb9\xe4?\xa18\x97Ag>\xaf(=Z4ն&<Vm\xd3/|\x1d\xa5yő]\xba\x1a\xa8\x84\xa6,V\xee\b\x1fr]\xc7\xdeM!\x9b8x\xb3u\x9e\xa3z\xc5nU\xf5!\xebV\x98\xebJb3\x82Ƣ\x81\xc9-\xf5\xeclT\xd2\x1bMѡ=\xf3i\x14\x11\xba\xa2\xd0%\xc1?J\x96\\\x9d4\xd2\xfdg\xbea\"S5\xab\xb6\xbc\x8b\xa2Aw9\x83\xbc\xb48\xb2K\xd35\xcc\xc8\xff\x1aR\x9a\x13h+\xb5E\x84\xba\x14R\x1e\xb1\xe0H\xbah\x17\xfb\x1a\x82\xb7 \x13\xb1N\xf3\xa5\xa4\xea\x0f\xb3\xd2r\xdb2\x12\xfe]\x15\x0e)+\xa1hd\b\x9b730\xc7iN!_\xa9\xa3\xc0\x1d\x05\xab\xbdc\xdel\x9c\xbak\x90\x0e\x95\xba\x87t\x9d\xdaSTwC\x96\x13C\x80.\xed\x81\x1e\xea\x81`c21Ҟ\x12'y\xa8*\x02@p=\xac0\xad<\xb1T\xe1\xc0\xfa)\v#q_\xc5\xecV\x92\xc6|Cu\xa3qQ\\G5B\xbe~\x83\"\xe1\xe5uZ\x97C\xef3\xb7\x80ö\a86\xed\x0e\xae\v\x18\x12\x0f)d\u07b6\xef\"\xee\xdf\xed?\xf8\x89\x90\xed\xc3n\xc8v\x1aZ\xe0\xb9r\xa6N\xe3\xdc\xf6e~\x03?\xa9\xc7\x05\xda}_\xcf\b\xd0J\xacV\x11;\xec\x8bC\x86\xee6vK\xe9#M\xae߄s\x93\xbf6_]\a\xbb鎎d\x06\xdd\xd1ep]\\\xb3z\xb2\xd9\u007f\xfc\xcef\x0fc5\xd5c\x01\xc8\x13\xc81%Ϛ\xc42Y\xec\xad\xefU\xf8,\x8a\xba\xacָ\x17Tm\a.\xab\xbe\x13\xd5oh7\xce\x17\xe7\x14\xb7Oc\xfd\xb2q[\xcf8!\xf9\xee;\x80\xe41$\xdf\xd8g\"$9\x87\xee\x1b6\t\xf7\x8bޔ</C\xfb+WŇ\x8d\x18i\x97/\v\xdb\x15\nG\xddIM\x95q:j}\xb8Ǎa\xfb\x96\xb3\xeb2\x8f)\\\x8b~#9\xa6|95\x1fCj\xcdτL\x98\xcfC\xee\xdb{!\x8d\x83\xbc\xae\xe7\xf6\x93\xa5u\xcd\xcee\xf0\x89Ő\x86zMlY\x17\xe5\x0f\xbc\x8d\x15p\xf9\xb8\x12b\x85\x89\xc5\xc3\x1b\x89N8\xf9ڤn\xa6\x8aG\xcdU\xfa OP\xe9_\x1aR\x97\xecTl\xec;e\x95qͯ\xed\xed\xa6;\xda\x10^\"\x92j\xediu\xcbڈ\x1f\xd7\\\xf1;\xba?\x04\x03\xbe\xadv\b0\x05\xb3H($\f\xff\x869kV\xf1l\x19q\xb9f\x98L\x00\xd4گXm\xd6+\xfd\xa7(L\xc04\x95\xec\n\xda\xc2\x12Go9#얚ENW\x8en,\xb3\xd7B\xd3\xed\x9da\x06ي\xcbkz\xed\xea\xe9\xa8mK#\x11(\xac?\xdbV,\xe8Y9\n\x82fޢX\x88\xdb\xd9\x15\xd4;2Ƀ\x9a\xe1`\xee4\u0530\xce\xfe\xf8\xf3\xd5\xe7\xab\u007f^\xc25\xf7\x82dq\xc0\xc0;,\xe8\t\x8c\x02\x16ZR'\x8b\xb9\x92\x0e\x99\x93\x9b[R)\x8c\xbb\xba\x8d3\xf3t\xfa4\xc0\xd6W=&l0\x1fP\xf9\x18\xfb7g\xb7\x83\xf1\xa0p\xe8`^|\x8e\a\xda`\x83\xb9\xfe\x19\x0f\x8c\xae\x83\xb9\xf9\x1d\x0f@\x97\xc1\x1c\xfe\xd3\xe0\x04}\xa4\xcf6p%\xf8\xa2u\xaf\x8b\xeciim1Y\x9ct\xebe\x81\x92L\xae]\xfc\xac\x18\f\xae=\u007f\ak\xbe\xf9\x01b\x96*\x95\xc2\x11\x0f\xb0i\xf5M\xe4_\x06\xff\xc0\xeeu\x82E\xaf*\a\x00'\t\x8b\x83Ϣ\xccneF\xba\n\xf5\x9b,>{\xe9\xb7\xdaX(\xb2d\xd0\x03\xb0\x98B;\x12\x10\x81\xa7F\x9a\xef\xfb5l\f\xadנ,\n\x1a\xe4tB\x9e29\xb6XHr\xe3\x9d\x18\xa1\xabЮ}\xc8\xc6\xffG`\xe4:p\xd7Èx\xf52^\xc6\xcbx\x19\xffW\xe3\xdf\x01\x00\x00\xff\xff\x05\x17n\xe1\x00(\x00\x00"))
	const prefix = "/assets/"
	manager = assets.NewManager(assetsLoader, prefix)
	App.SetAssetsManager(manager)
	assetsHandler := assets.Handler(manager)
	App.Handle("^"+prefix, func(ctx *app.Context) { assetsHandler(ctx, ctx.R) })
	App.AddTemplateVars(map[string]interface{}{
		"AdSense":   func() interface{} { return AdSense },
		"Chitika":   func() interface{} { return Chitika },
		"Top":       Top,
		"Bottom":    Bottom,
		"Fixed":     Fixed,
		"Sizes":     func() interface{} { return sizes },
		"providers": func() interface{} { return providers },
	})
	templatesLoader := loaders.TgzLoader(loaders.RawString("\x1f\x8b\b\x00\x00\tn\x88\x02\xff\xecX\xdfo\xfa6\x10癿\xe2\xe4=\xac|\xf5%\x10J[\xa9\x05\xd4jӤ=L\x9aV\xedi\x9d*\a\x1b\xe25\xc4Q\xce)E\xd0\xff}g'@\xf81\xb5\xebJ\xbb~\x9b{\x80$w\xe7;\xdb\xf7\xe3cs\x812F\xe9\x85f\x12\xd5\x0eD\xed\x9c\xfe\xe9\xbf\xe3\xb7\xfd\xe5s\xfe\xdd\xf7\xbb\xbe_\x83\xda[P\x86\x86\xa7d\xfe\xb5'\xf9Ah>o}\xa9\x03\xc8\a#c\x81\xe7\x10\xf0\"\x18\xea_Z\x8f\x8f\xf5\xfa|\x0eB\x8eT,\x81\xdd\xde\x16\xc1B\xff\f\x88I<5\x02\x1e\v\xb8\xbc\x12\xd7.\x8c~͂Ha(S\xf0\xac\x04@O\xa8{\x18F\x1c\xb1\xcfH\xbdY\f\x01\xf69\xd0\x0f\xee\xdf\xe8\xc4>\xe7\xc3]^\x87\\\xe8)i;\x1e\xba7b\x91{\xf4\x8d\rhL\x00z7r\x92Dܐ_V,TB6\x83\xcc\x18\x1d\xb3ܰ\xa5\x9e\x8a\xb1d<\x98\x8d\xb5\x1eG+\xe3\xac\x10\x03@3\x8bd\x9f\t\x854\xe6\xec<\x88\xf4\xf0n\xcd\x15\xdcpr\xbc9\x8c\x94\x8cM\x9f\x91\xf5=\x13&\xe7v40ҹ\xbc\xb7\x97=\xd2鄓\x00ό.fFN\xb7\xc8\xeb\xd5\v\x0eS\x95\x98\xc1Q\xd9\xfd>LUL\xab\xe2\x95?.\x16\xf0ǟ\r/\xc90<\x9a?6.z\xadB\xb5\x9e\x0fJ\xfb0\xb0;&#Z}\xbb\xce\xfbfP߳\xb42Mu\xca\xe0\xc8\f\x81]\td\xc0~Q\x88*\x1eC\xb2\xd2\xfc\xf9ǯ\x80\xd2\xd0\xc2z\xbb\xc3\x1a\r3\x9d\xa5\x1b\xe2\xacQ\x04\x90s\xe7\xdf\x1a\xb6\xcb\xeal&\xb4\xb5\xa0\xc8.\x82\t\xe5Z\x9f\xa7\xe3lB{\x05M\x90\xdexc\xe8\x1b;֢\xf0\xb2\xf9\x9b\xc4DǨ\xee\xe5\r#\x96\xdf9\ue79c\xde\xd8\x10b\xb0r\xd1E^驜\x14\xbb\x03\xb1\xfd\x81\x9f\xae\x05\xf6\xc4\xf0fn-sg\xbdi\xcf1\xdc\xfcI=H\xf1\xa4y\xd8|m\x8e\x9c\u058b}\xaaU\xf4\x91\x89\x02D\x1a<h\xfb\u007f\xaa\xff\xb7\xfdvg\xbb\xffw\xce\xdaU\xff\u007f\xc3\xfe\xef\x1a .\x82,\x16<\x88\xe4\xb9-\x12\xde\x10Ѳ\\\x1b)\xf3l3\xf2\xfe\xc2F\x0e\x11\xaa\x1c\xfa\xc84\f\x95Qw\xfc=\xf1?\xa5\xff\xc9\x0e\xfeo\x9fU\xf9\xff?\xc0\xff\x1b\U0003f215%\xfc/@\x86\x03\xb3+l\x97\xa3\xe3\x1f\x8a\xa8\xda@ǹ$*\xe1d\xb4=\"\xb0\x1ci\x95pJa\xe3\xf0\xe7\x83\x1d S\x9aj\xe1\xfe\u007f\au\xe5\x15\xf3^b\xf7\xf51\xdds\\\xaa\xca\xe2\xa7!w\xc8:,\xfc{\x06\xfe;ݮ\xfft\x10\xab\xea\xff\xdb\xd4\xffU\xf5)\x9d\xb9\U000f2414\x8bM\xce\x19\x14w\x19\xbdVR\x15\x8co\x81J\x8d\xf1`U\xe0\xa9\xfc\xef\x1ew\xb7\xf2\xbf\xd3\xe9T\xe7\xbf\xf7\xc8\xff}\xf7\xa8\x16\"\xc2\xef(ay=\x1ak\x12WqDZ\xe2+\xa0\x06e\xbeG U!c\x98\x86\xf4\x93\xb9k:.\xdcE*8$\xe9n!y\xb9\xa2l\xd8ں\x81\xb5&.\x18\x84\xa9\x1c\xf5\xd9wl\xd0Ä\xc7\xdbʮ\x1a\xado\a\xdd'W\x9b\xac\xf0\xae\n\x86z\xba\xa5\xe2>\x95TZ\xbc\xaaj\x15UT\xd1g\xa1\xbf\x03\x00\x00\xff\xffģ\x9a\x91\x00\x1e\x00\x00"))
	App.SetTemplatesLoader(templatesLoader)
	tmpl_assets_html := template.New(templatesLoader, manager)
	tmpl_assets_html.Funcs(map[string]interface{}{
		"t":       func(_ ...interface{}) interface{} { return nil },
		"tn":      func(_ ...interface{}) interface{} { return nil },
		"tc":      func(_ ...interface{}) interface{} { return nil },
		"tnc":     func(_ ...interface{}) interface{} { return nil },
		"reverse": func(_ ...interface{}) interface{} { return nil },
	})
	if err := tmpl_assets_html.Parse("assets.html"); err != nil {
		panic(err)
	}
	App.AddHook(&template.Hook{Template: tmpl_assets_html, Position: assets.None})
}
