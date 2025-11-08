# EzPlay
---

I watch lecture videos stored on my local device and typically play them using `mpv`. However, manually typing the lecture name and subtitle file path for the `mpv` command became quite a hurdle. So I decided to automate it. I built this simple tool using Go and [Fyne](https://fyne.io/) that features a search menu and an option to include subtitle files or play without them.

## Images
<p align="center">
    <img src="screenshots/AppStartupDemo.png" alt="AppStartupDemo - EzPlay" width="45%" />
    <img src="screenshots/LecturePlayDemo.png" alt="LecturePlayDemo - EzPlay" width="45%" />
</p>


## Installation

1. Clone the repository:
```bash
git clone https://github.com/kshrs/ezplay.git
cd ezplay
```

2. Build the executable:
```bash
go mod tidy
mkdir build
go build -o build/ezplay
```

3. Move the executable to your PATH:
```bash
sudo mv build/ezplay /usr/local/bin/
```

Now you can run `ezplay` from anywhere in your terminal.

---

**Tested on Arch Linux**

