# voicechannel-discordbot

Cooler Bot, Erstellt beim beitretten eines besonderen Channels eine neue Kategorie mit Voice und Textchannel. Auf den Textchannel haben nur leute Zugriff, die auch aktuell im Voicechannel sind. Nachdem der letzte gegangen ist, wird die Kategorie mit den Voice und Textchannel wieder gelöscht ... Also temporäre Channel

Gibt nen paar Gimmicks, wie zum Beispiel spezielle Channelnamen, die man per Command hinzufügen kann (oder per editieren der sql-lite DB

## Setup

1. Eine Config erstellen. Die sollte unter config.json im Hauptverzeichnis des Programms liegen und folgendes beinhalten:
```
{
    "Token": "",
    "MasterChannel": "",
    "KategoriId": "",
    "Everybody": ""
}
```
Wobei Token das Discord-Developer-Bot-Token ist, was man sich unter Discord-Dev generieren kann,
Masterchannel die ID des Channels, bei dem die Temporären Channel erstellt werden sollen
Kategorie... kA, ist schon lange her. Glaube das ist veraltet und sollte ich mal entfernen
Everybody - Die ID von @Everybody.. die ist Serverspezifisch.. Ich denke da kann ich bestimmt irgendwann noch was fixen, damit das klappt ohne

2. `Docker-compose up --build -d`
3. Den Bot zum Server einladen und admin Permissions geben (todo: Fixen, die Permissions braucht er nur, weil ich noch zu faul war den textchannel erstmal öffentlich zu erstellen, den Bot die Permissions zu geben und ihn dann erst unsichtbar für alle zu machen... Ist alles recht nervig. Der fix war an der Stelle einfach Admin Permissions.. Auch wenn das keine gute Lösung ist. Ich hatte keine Lust mehr 🤷‍♂️
4. Sollte klappen. Wenn nicht, mach nen Issue auf.. oder nicht. Ich denke das Ding werde soweiso nur ich benutzen .. 
