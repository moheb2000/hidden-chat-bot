// Code generated by running "go generate" in golang.org/x/text. DO NOT EDIT.

package translations

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
)

type dictionary struct {
	index []uint32
	data  string
}

func (d *dictionary) Lookup(key string) (data string, ok bool) {
	p, ok := messageKeyToIndex[key]
	if !ok {
		return "", false
	}
	start, end := d.index[p], d.index[p+1]
	if start == end {
		return "", false
	}
	return d.data[start:end], true
}

func init() {
	dict := map[string]catalog.Dictionary{
		"en_US": &dictionary{index: en_USIndex, data: en_USData},
		"fa_IR": &dictionary{index: fa_IRIndex, data: fa_IRData},
	}
	fallback := language.MustParse("en-US")
	cat, err := catalog.NewFromMap(dict, catalog.Fallback(fallback))
	if err != nil {
		panic(err)
	}
	message.DefaultCatalog = cat
}

var messageKeyToIndex = map[string]int{
	"Are you sure you want change the hidden link? Your previous link will be invalid! 🚨🚫": 9,
	"Audio":                   16,
	"Banned succesfully! 🔴✅":  33,
	"Blocked succesfully! 🔒✅": 36,
	"Cancel ❌":                8,
	"Document":                17,
	"Gif":                     12,
	"Photo":                   13,
	"Send your message ✏️:":   30,
	"Sticker":                 11,
	"Text":                    10,
	"The user you trying to send this message to, changes the hidden link. Maybe you could ask the user to get the new one! ⚠️": 22,
	"There is a problem in our servers. Please be patient and try later! ⚠️":                                                    38,
	"This link is not valid. Maybe you need to contact somehow to the link's owner and tell this problem! ⚠️":                   29,
	"Unbanned successfully! 🟢✅":  34,
	"Unblocked successfully! 🔓✅": 37,
	"Video":                      14,
	"Voice":                      15,
	"Yes, I'm sure ✅":            7,
	"You are banned by bot's admin and you can't use bot anymore! ⛔❌🔴":            39,
	"You are not in sending mode... ⛔":                                            21,
	"You blocked by the user you're trying to send the message! 🔒😿":               23,
	"Your message sended successfully! 📨😍":                                        28,
	"Your message type is limited by reciever or isn't supported by this bot. 🔒🥹": 27,
	"about_message":                         20,
	"message_restrictions_settings_message": 19,
	"settings_message":                      6,
	"start_message":                         3,
	"ℹ️ About":                              0,
	"⚙️ Settings":                           1,
	"⬅️ Back":                               18,
	"💬 Reply":                               24,
	"🔒 Block":                               25,
	"🔓 Unblock":                             35,
	"🔗 Change Hidden Link":                  4,
	"🔗 Get Hidden Link":                     2,
	"🔴 Ban":                                 31,
	"🚨 Report":                              26,
	"🚫 Message Restrictions":                5,
	"🟢 Unban":                               32,
}

var en_USIndex = []uint32{ // 41 elements
	// Entry 0 - 1F
	0x00000000, 0x0000000d, 0x0000001d, 0x00000032,
	0x0000014b, 0x00000163, 0x0000017d, 0x00000214,
	0x00000226, 0x00000231, 0x0000028c, 0x00000291,
	0x00000299, 0x0000029d, 0x000002a3, 0x000002a9,
	0x000002af, 0x000002b5, 0x000002be, 0x000002ca,
	0x00000319, 0x00000420, 0x00000443, 0x000004c1,
	0x00000505, 0x00000510, 0x0000051b, 0x00000527,
	0x00000579, 0x000005a4, 0x00000610, 0x0000062a,
	// Entry 20 - 3F
	0x00000633, 0x0000063e, 0x0000065a, 0x00000679,
	0x00000686, 0x000006a3, 0x000006c3, 0x0000070e,
	0x00000756,
} // Size: 188 bytes

const en_USData string = "" + // Size: 1878 bytes
	"\x02ℹ️ About\x02⚙️ Settings\x02🔗 Get Hidden Link\x02Hello 👋, Welcome to " +
	"*%[1]v*!🥰\x0a\x0aYou can create hidden chat links and others can send yo" +
	"u messages without knowing your username😻😻.\x0a\x0aJust click on *\x22🔗 " +
	"Get Hidden Link\x22* button and paste your link in your social media to " +
	"start recieving wonderful messages👻🤡🙈!\x02🔗 Change Hidden Link\x02🚫 Mess" +
	"age Restrictions\x02Here is your account control center! 🛠️🪛🙆\x0a\x0aYou" +
	" can change your hidden link 🔗, set allowed/disallowed message types ✅❌ " +
	"and more! 😇\x02Yes, I'm sure ✅\x02Cancel ❌\x02Are you sure you want chan" +
	"ge the hidden link? Your previous link will be invalid! 🚨🚫\x02Text\x02St" +
	"icker\x02Gif\x02Photo\x02Video\x02Voice\x02Audio\x02Document\x02⬅️ Back" +
	"\x02You can set what type of messages you want to recieve or not from ot" +
	"hers! 💯\x02Create *hidden chat links* to send messages without revealing" +
	" your username. 🔒\x0aShare your link, stay anonymous, and enjoy private " +
	"conversations. 💬\x0aYour identity remains secure while you chat freely. " +
	"👻\x0a\x0aLet’s keep your chats secret and fun! 🤖✨🎉\x02You are not in " +
	"sending mode... ⛔\x02The user you trying to send this message to, change" +
	"s the hidden link. Maybe you could ask the user to get the new one! ⚠️" +
	"\x02You blocked by the user you're trying to send the message! 🔒😿\x02💬 R" +
	"eply\x02🔒 Block\x02🚨 Report\x02Your message type is limited by reciever " +
	"or isn't supported by this bot. 🔒🥹\x02Your message sended successfully! " +
	"📨😍\x02This link is not valid. Maybe you need to contact somehow to th" +
	"e link's owner and tell this problem! ⚠️\x02Send your message ✏️:\x02🔴 B" +
	"an\x02🟢 Unban\x02Banned succesfully! 🔴✅\x02Unbanned successfully! 🟢✅\x02" +
	"🔓 Unblock\x02Blocked succesfully! 🔒✅\x02Unblocked successfully! 🔓✅" +
	"\x02There is a problem in our servers. Please be patient and try later! " +
	"⚠️\x02You are banned by bot's admin and you can't use bot anymore! ⛔❌🔴"

var fa_IRIndex = []uint32{ // 41 elements
	// Entry 0 - 1F
	0x00000000, 0x00000019, 0x0000002f, 0x00000053,
	0x0000022b, 0x0000024d, 0x00000271, 0x00000364,
	0x0000037e, 0x00000389, 0x00000423, 0x0000042a,
	0x00000437, 0x0000043e, 0x00000445, 0x0000044e,
	0x00000455, 0x0000045c, 0x00000465, 0x00000479,
	0x00000505, 0x00000686, 0x000006ce, 0x00000795,
	0x000007ef, 0x000007fd, 0x0000080b, 0x0000081b,
	0x000008d6, 0x0000090d, 0x00000995, 0x000009b5,
	// Entry 20 - 3F
	0x000009c5, 0x000009e0, 0x00000a0b, 0x00000a48,
	0x00000a5a, 0x00000a83, 0x00000ab0, 0x00000b3a,
	0x00000bc9,
} // Size: 188 bytes

const fa_IRData string = "" + // Size: 3017 bytes
	"\x02ℹ️ درباره ما\x02⚙️ تنظیمات\x02🔗 دریافت لینک پیام\x02سلام 👋، به *%[1]" +
	"v* خوش اومدی!🥰\x0a\x0aاینجا میتونی لینک چت ناشناس بسازی و بقیه بدون اینک" +
	"ه یوزرنیم تو رو بدونن بهت پیام بدن😻😻.\x0a\x0aفقط کافیه روی دکمه *«🔗 دری" +
	"افت لینک پیام»* کلیک کنی و لینکتو بگذاری توی شبکه های اجتماعی تا شروع ک" +
	"نی به دریافت پیام های شگفت انگیز بقیه👻🤡🙈!\x02🔗 تغییر لینک پیام\x02🚫 محد" +
	"ودیت های پیام\x02اینجا تنظیمات اکانت توئه! 🛠️🪛🙆\x0a\x0a میتونی لینک پیا" +
	"متو تغییر بدی 🔗، مشخص کنی چه پیام هایی برات بیان یا نیان ✅❌ و خیلی چیزا" +
	"ی بیشتر! 😇\x02آره، مطمئنم ✅\x02لغو ❌\x02مطمئنی میخوای لینک پیامتو تغییر" +
	" بدی؟ اگه اینکارو کنی لینک قبلی دیگه معتبر نیست! 🚨🚫\x02متن\x02استیکر\x02" +
	"گیف\x02عکس\x02فیلم\x02ویس\x02صدا\x02فایل\x02➡️ بازگشت\x02اینجا میتونی م" +
	"شخص کنی چه پیام هایی دوست داری که از بقیه دریافت کنی یا نکنی! 💯\x02*لین" +
	"ک چت ناشناس* بساز و بدون اینکه یوزرنیمت رو به بقیه نشون بدی چت کن. 🔒" +
	"\x0aلینکتو منتشر کن، مخفی بمون و از چت ناشناس لذت ببر. 💬\x0aدر حالی که ر" +
	"احت چت می کنی، اینکه کی هستی مخفی می مونه. 👻\x0a\x0aبیا چتاتو امن و جال" +
	"ب کنیم! 🤖✨🎉\x02شما در وضعیت ارسال پیام قرار ندارید... ⛔\x02کاربری که می" +
	"خوای بهش پیام بفرستی، لینک پیامشو عوض کرده. شاید بتونی دوباره ازش بخوای" +
	" تا لینک جدیدو بهت بده! ⚠️\x02کاربری که میخوای بهش پیام بفرستی بلاکت کرد" +
	"ه! 🔒😿\x02💬 پاسخ\x02🔒 بلاک\x02🚨 گزارش\x02دریافت کننده پیام نوع پیامی که " +
	"میخوای بهش بفرستی رو محدود کرده یا این نوع توسط بات پشتیبانی نمیشه. 🔒🥹" +
	"\x02پیامت با موفقیت ارسال شد! 📨😍\x02این لینک معتبر نیست. شاید لازم باشه " +
	"یه جوری به صاحب لینک، این مشکلو بگی! ⚠️\x02پیامتو بفرست ✏️:\x02🔴 مسدود" +
	"\x02🟢 رفع مسدودیت\x02با موفقیت مسدود شد! 🔴✅\x02با موفقیت از مسدودیت خارج" +
	" شد! 🟢✅\x02🔓 آنبلاک\x02با موفقیت بلاک شد! 🔒✅\x02با موفقیت آنبلاک شد! 🔓✅" +
	"\x02یه مشکلی توی سرورها به وجود اومده. لطفا صبور باش و بعدا دوباره امتحا" +
	"ن کن! ⚠️\x02ادمین بات شما رو مسدود کرده و نمی تونید از این به بعد از با" +
	"ت استفاده کنید! ⛔❌🔴"

	// Total table size 5271 bytes (5KiB); checksum: 4215A8D5
