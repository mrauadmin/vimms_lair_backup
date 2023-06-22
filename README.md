# Vimm's Lair (The Vault) Archive
___
- [Description](#description)

- [Instalation](#instalation)

- [License](#license)
___

### Description:

While [vim.net](https://vimm.net) is an amazing repository of over 60k+ games, it is a one website run by one (!) dude with their own servers. It could go underwater tomorrow from a cease and desist from some random corporation and they would not have the resources to move, deploy a mirror in time or defend from the lawsuit. That's why this script exists. It downloads and then backups the entire collection to WebArchive (WIP), Torrents (WIP) and IPFS (WIP).

### Instalation:

1. Clone the code to folder of your choice:
```
$ git clone https://github.com/mrauadmin/vimms_lair_backup
```

2. While this code is still in development you need to edit the file itself to change the output folder:

**Linux Native:**
```
$ cd vimms_lair_backup
$ nano main.go 
```

Find the `const Path_to_save = "Z:\gry"` (around the line 32) and change it to the folder of your choosing. Click `Ctrl + O` and then `Ctrl + X`.

**Windows:**

Open the `vimms_lair_backup` and open the `main.go` file and find the `const Path_to_save = "Z:\gry"` (around the line 32) and change it to the folder of your choosing.

3. Run the code:

**Linux Native:**
```
$ cd vimms_lair_backup
$ go run main.go
```

**Windows:**

Go to the `vimms_lair_backup` and then run this command in your editor or in cmd: 
```
$ go run main.go
```

4. Happy Downloading!

### License:

[cc-by-nc-sa]: http://creativecommons.org/licenses/by-nc-sa/4.0/
[cc-by-nc-sa-image]: https://licensebuttons.net/l/by-nc-sa/4.0/88x31.png
[cc-by-nc-sa-shield]: https://img.shields.io/badge/License-CC%20BY--NC--SA%204.0-lightgrey.svg

This work is licensed under a
[Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License][cc-by-nc-sa].

[![CC BY-NC-SA 4.0][cc-by-nc-sa-shield]][cc-by-nc-sa]

[![CC BY-NC-SA 4.0][cc-by-nc-sa-image]][cc-by-nc-sa]
