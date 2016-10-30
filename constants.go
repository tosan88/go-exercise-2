package main

//https://www.alien.net.au/irc/irc2numerics.html
const (
	RPL_WELCOME              = "001"
	RPL_ENDOFMOTD            = "376"
	PING                     = "PING"
	JOIN                     = "JOIN"
	MESSAGE_CMD              = "PRIVMSG"
	RPL_NAMREPLY             = "353"
	KICK                     = "KICK"
	PART                     = "PART"
	QUIT                     = "QUIT"
	ERR_ERRONEUSNICKNAME     = "ERR_ERRONEUSNICKNAME"
	ERR_ERRONEUSNICKNAME_NUM = "432"
	ERR_NICKNAMEINUSE        = "ERR_NICKNAMEINUSE"
	ERR_NICKNAMEINUSE_NUM    = "433"
	ERR_NICKCOLLISION        = "ERR_NICKCOLLISION"
	ERR_NICKCOLLISION_NUM    = "436"
)

var randomText = []string{
	`Cats can drink salt water to survive: their kidneys can filter out salt!`,
	`Cats sleep for about 70% of their lives.`,
	`Disneyland owns over 200 cats.`,
	`Cats usually reserve meows in order to communicate with humans.`,
	`The technical term for a cat's hairball is a "bezoar".`,
	`A group of cats is called a "clowder".`,
	`It takes approximately 24 cat skins to make a coat.`,
	`There are about 40 different breeds of cats.`,
	`Smuggling a cat out of ancient Egypt was punishable by death.`,
	`Cats are great at detecting hidden microphones planted by Russian spies!`,
	`All cats sheath their claws at rest except for the cheetah.`,
	`A cat lover is called an Ailurophile.`,
	`Cats sweat through their paws.`,
	`Cats spend 30% of their waking hours on cleaning themselves.`,
	`The first cat in space was a French cat named Felicette.`,
	`Around the world, cats take a break to nap —a catnap— 425 million times a day.`,
	`A female cat is also known to be called a "queen" or a "molly."`,
	`The Snow Leopard, a variety of the California Spangled Cat, always has blue eyes.`,
	`Each side of a cat's face has about 12 whiskers.`,
	`Today, cats are living twice as long as they did just 50 years ago.`,
	`Cats are unable to detect sweetness in anything they taste.`,
	`Twenty-five percent of cat owners use a blow drier on their cats after bathing.`,
	`Cats have over 100 sounds in their vocal repertoire, while dogs only have 10.`,
	`A third of cats' time spent awake is usually spent cleaning themselves.`,
	`Unlike most other cats, the Turkish Van breed has a water-resistant coat and enjoys being in water.`,
	`Some cats can survive falls from as high up as 65 feet or more.`,
	`Cats have a 5 toes on their front paws and 4 on each back paw.`,
	`Cats have 24 more bones than humans.`,
	`Collectively, kittens yawn about 200 million time per hour.`,
	`Cats have a strong aversion to anything citrus.`,
	`The Maine Coon is appropriately the official State cat of its namesake state.`,
	`A fingerprint is to a human as a noseprint is to a cat.`,
	`Sir Isaac Newton, among his many achievements, invented the cat "flap" door.`,
	`The two outer layers of a cat's hair are called, respectively, the guard hair and the awn hair.`,
	`Cats greet one another by rubbing their noses together.`,
	`Most cats don't have eyelashes.`,
	`Cats invented The Internet.`,
	`Cats are 110% better than dogs.`,
	`Every cat's nose is unique, and no two nose-prints are the same.`,
	`Cats rub against people to mark them as their territory.`,
	`When cats are happy or pleased, they momentarily squeeze their eyes shut.`,
}
