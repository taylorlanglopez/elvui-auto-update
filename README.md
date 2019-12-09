# elvui-auto-update
A script that you can run manually after you notice the message 'your elvui is out of date' 
requires a manual edit to insert your retail interface path

This legit just runs some go code to download the link from tukui.org and then save the .zip to the interface folder,
then extracts the zip and copies it to the AddOns folder, very limited testing, I suggest a backup if you actually use this.

This took like 10 minutes to make so use caution, I do not guarantee the safety of this thing or even downloading the correct file.
The program only works in my limited testing because it happens to be the only .zip file that is included in the response to the GET request
at the url -> https://tukui.org/download.php?ui=elvui

If they ever update the page to include more than one .zip, this will return the first .zip link found and then download whatever that is,
this tool is for my personal convenience
