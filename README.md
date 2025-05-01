# 🛠️ Go Multi-Threaded File Downloader

This is a simple, concurrent file downloader written in **Go**, which splits a file into multiple byte-range sections and downloads them in parallel. Ideal for learning how to work with HTTP ranges, goroutines, and basic I/O in Go.

---

## 🚀 Features

- Downloads files using HTTP **range requests**
- Splits the file into **multiple sections**
- Downloads each section in a **separate goroutine**
- Merges all the parts into a final file
- Cleans up temporary section files
- Shows file size and progress logs

---

## 📦 Requirements

- Go 1.18 or higher
- Internet connection 😄

---

## 🔧 How It Works

1. Sends a `HEAD` request to get the file size
2. Calculates byte ranges based on the number of sections
3. Starts parallel downloads using goroutines for each byte range
4. Writes each section to a temporary file
5. Merges all temporary files into one final file
6. Deletes the temporary section files

---

## 🧪 Example

```bash
go run main.go
```

It will download this sample MP4 file:

```
https://file-examples-com.github.io/uploads/2017/04/file_example_MP4_1920_18MG.mp4
```

Into a file called `final.mp4` in the project directory.

---

## 📝 Output Example

```bash
Making Connections from DO Func...
Got 200
Size of this file is: 18329232
Size of each section : 1832923
Section after shaped:  [[0 1832923] [1832924 3665847] ... [16496308 18329231]]
Downloaded 1832924 bytes for section 0 [0 1832923]
...
1832924 bytes merged
...
Downloade completed in 12.31 seconds
```

---

## 📁 Files

- `main.go`: Core implementation
- `section-<n>.tmp`: Temporary section files (auto-deleted after merge)
- `final.mp4`: The final merged output

---

## ⚠️ Notes

- Ensure the server supports **HTTP range requests** (most do).
- The total number of sections can be tuned with the `TotalSection` field.
- For very large files, increase concurrency with caution — high parallelism can overwhelm weaker servers.

---

## 🧹 TODO

- Progress bar
- Retry logic on failed downloads
- CLI interface for URL & section count

---

## 📜 License

MIT — use it for anything, just give credit if it helps!
