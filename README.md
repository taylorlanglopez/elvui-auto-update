# elvui-auto-update
A go program you can run from inside VSCode or using a WSL terminal to update your elvui automatically. 

To work properly, the code must be edited with your retail interface path so that the extraction works correctly.

This legit just runs some go code to download the link from tukui.org and then save the .zip to the interface folder,
then extracts the zip and copies it to the AddOns folder, very limited testing, I suggest a backup if you actually use this.

This will only download and extract the specific file that is patterned as elvui-##.##.zip, so if they add more .zip's to the page this will still only download the new update.

Requires GO to be installed, and preferably WSL so that you could run the .sh script from a terminal and it will just execute the GO code. 

If you modify the .sh file with the path of wherever your downloaded the repo to, you can place the .sh file on your desktop and exectue it quickly. This tool is mainly for my convieniece since I just use a MinGW64/gitBash term to execute it from my /home/ whenever I get the in game message "your elvui is out of date".
