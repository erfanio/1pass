# 1pass
GUI for [1password cli](https://support.1password.com/command-line-getting-started/)  
Since linux doesn't have a 1password client, this is a simple way to access passwords in 1password

## Status
This is by no mean a full featured 1password client, this is just a convenient way to access your 1password vault (potentially from a key binding).
This project is still work in progress, but essential parts are working now.

![Demo](https://thumbs.gfycat.com/SizzlingHonoredHamadryas-max-14mb.gif)

After logging in, you're going to be remembered for 30min.  
From the list you can press `ctrl-c` on an item to copy the password into your clipboard (if there is a password).  
Pressing enter will open the item in details.

# Build
You need to install [qt](https://github.com/therecipe/qt/wiki/Installation) for go first.
You can find [building instructions](https://github.com/therecipe/qt/wiki/Getting-Started#starting-application) for qt applications too.
Since custom constructor/signal is used, qtmoc needs to run before qtdeploy.

