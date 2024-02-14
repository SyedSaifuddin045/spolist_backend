package song

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type SongDownloadRequest struct {
	SongID   string `json:"songID"`
	SongLink string `json:"songLink"`
}

func HandleSongDownload(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		SendSong(w, r)
	case http.MethodPost:
		StartSongDownload(w, r)
	default:
		http.Error(w, "Method Not allowed", http.StatusMethodNotAllowed)
	}
}

func StartSongDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not allowed", http.StatusMethodNotAllowed)
		return
	}

	var songRequest SongDownloadRequest
	// Decode the JSON request body into SongDownloadRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&songRequest)
	if err != nil {
		http.Error(w, "Error decoding JSON request", http.StatusBadRequest)
		return
	}

	// Now you can access the values of SongID and SongLink
	fmt.Println("Song ID:", songRequest.SongID)
	fmt.Println("Song Link:", songRequest.SongLink)

	songPath := fmt.Sprintf("static/%s.mp3", songRequest.SongID)
	songAlreadyExists := isFilePathValid(songPath)
	if songAlreadyExists {
		w.WriteHeader(http.StatusOK)
		return
	}

	command := fmt.Sprintf("spotify_dl -l %s -o %s -mc 4", songRequest.SongLink, songRequest.SongID)
	cmd := exec.Command("/bin/bash", "-c", command) // Use "/bin/bash" to execute the command

	// Create pipes for capturing standard output and standard error
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating StdoutPipe: %v", err), http.StatusInternalServerError)
		return
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating StderrPipe: %v", err), http.StatusInternalServerError)
		return
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error starting command: %v", err), http.StatusInternalServerError)
		return
	}

	// Create goroutines to read and display the output in real-time
	go displayOutput(stdout)
	go displayOutput(stderr)

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error waiting for command: %v", err), http.StatusInternalServerError)
	}

	id := songRequest.SongID
	mp3Path, err := findMP3Path(id)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(mp3Path)

	destinationPath := "static/"

	err = moveMP3File(mp3Path, destinationPath, id)
	if err != nil {
		fmt.Println(err)
		return
	}
	removeError := removePath(id)
	if removeError != nil {
		fmt.Printf("Error deleting folder: %v\n", removeError)
	} else {
		fmt.Printf("Folder %s and its contents deleted successfully\n", id)
	}
}

func findMP3Path(rootPath string) (string, error) {
	var mp3Path string

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file is an MP3 file
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".mp3") {
			mp3Path = path
			return filepath.SkipDir // Stop walking after the first MP3 file is found
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if mp3Path == "" {
		return "", fmt.Errorf("no MP3 file found in the given directory and its subdirectories")
	}

	return mp3Path, nil
}

func isFilePathValid(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func SendSong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not allowed", http.StatusMethodNotAllowed)
		return
	}

	// You may want to validate the songID parameter in the URL or perform any other necessary checks

	songID := r.URL.Query().Get("songID")

	if songID == "" {
		http.Error(w, "Missing songID parameter", http.StatusBadRequest)
		return
	}

	songFileName := fmt.Sprintf("static/%s.mp3", songID)

	// Open the downloaded song file
	file, err := os.Open(songFileName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening song file: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set the response headers
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", songFileName))

	// Copy the file content to the response writer
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending song file: %v", err), http.StatusInternalServerError)
		return
	}
}

func moveMP3File(sourcePath, destinationFolder, newFilename string) error {
	// Open the source file
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination folder if it doesn't exist
	err = os.MkdirAll(destinationFolder, os.ModePerm)
	if err != nil {
		return err
	}

	// Create the destination file with the new filename
	destinationPath := filepath.Join(destinationFolder, newFilename+".mp3")
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Close the destination file before removing the source file
	err = destinationFile.Close()
	if err != nil {
		return err
	}

	// os.Remove(sourcePath)

	fmt.Printf("Moved %s to %s\n", sourcePath, destinationPath)

	return nil
}

func removePath(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Printf("Error removing path %s: %v\n", path, err)
	}
	return err
}

// displayOutput reads from the given reader and displays the output in real-time
func displayOutput(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		// You may want to send this output to the client or log it
		fmt.Println(scanner.Text())
	}
}
