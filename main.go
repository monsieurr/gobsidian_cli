package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Structure de configuration
type Config struct {
	VaultPath string `json:"vaultPath"`
}

// Renvoie le chemin du fichier de configuration dans le répertoire personnel
func getConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		// Si impossible d'obtenir le home, on le place dans le répertoire courant
		return "./.vaultconfig.json"
	}
	return filepath.Join(home, ".vaultconfig.json")
}

// Charge la configuration depuis le fichier ou crée une config par défaut si le fichier n'existe pas
func loadConfig() (Config, error) {
	configFilePath := getConfigFilePath()
	var config Config
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// Chemin par défaut si aucune config n'existe
		// config = Config{VaultPath: "/Users/mac/Documents/"}
		// Sauvegarde la config par défaut pour les prochaines utilisations
		err = saveConfig(config)
		return config, err
	}
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
}

// Sauvegarde la configuration dans le fichier de config
func saveConfig(config Config) error {
	configFilePath := getConfigFilePath()
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFilePath, data, 0644)
}

var vaultPath string

// Chargement de la configuration au démarrage
func init() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Erreur lors du chargement de la configuration :", err)
		// Si nécessaire de définir un path par défaut en utilisant pas le fichier vaultconfig.json
		// vaultPath = "/Users/mac/Documents"
	} else {
		vaultPath = config.VaultPath
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("-- Bienvenue dans Obsidian sur Terminal. --")
	fmt.Println("Tapez 'help' pour afficher la liste des commandes ou 'quit' pour quitter.")

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}
		command := args[0]
		switch command {
		case "quit", "exit":
			fmt.Println("Au revoir !")
			return
		case "help":
			printHelp()
		case "open":
			openObsidian()
		case "new":
			if len(args) < 2 {
				fmt.Println("Usage: new <nom_de_la_note>")
				continue
			}
			createNote(args[1])
		case "write":
			if len(args) < 2 {
				fmt.Println("Usage: write <nom_de_la_note>")
				continue
			}
			writeNote(args[1])
		case "delete":
			if len(args) < 2 {
				fmt.Println("Usage: delete <nom_de_la_note>")
				continue
			}
			deleteNote(args[1])
		case "list":
			listNotes()
		case "search":
			if len(args) < 2 {
				fmt.Println("Usage: search <mot_clé>")
				continue
			}
			searchNote(args[1])
		case "vault":
			fmt.Println("Le coffre se trouve à :", vaultPath)
		case "setvault":
			if len(args) < 2 {
				fmt.Println("Usage: setvault <chemin_du_coffre>")
				continue
			}
			updateVaultConfig(args[1])
		case "push":
			pushChanges()
		default:
			fmt.Println("Commande inconnue :", command)
		}
	}
}

func printHelp() {
	fmt.Println("Commandes disponibles :")
	fmt.Println("  open              : Lancer Obsidian")
	fmt.Println("  new <nom>         : Créer une nouvelle note (vérifie si elle existe déjà)")
	fmt.Println("  write <nom>       : Éditer une note existante")
	fmt.Println("  delete <nom>      : Supprimer une note existante")
	fmt.Println("  list              : Afficher toutes les notes")
	fmt.Println("  search <mot_clé>  : Rechercher des notes contenant le mot-clé dans leur nom")
	fmt.Println("  vault             : Afficher l'emplacement actuel du coffre")
	fmt.Println("  setvault <chemin> : Changer et sauvegarder le chemin du coffre")
	fmt.Println("  push              : Envoyer les modifications en ligne (Git push)")
	fmt.Println("  help              : Afficher cette aide")
	fmt.Println("  quit ou exit      : Quitter le CLI")
}

func openObsidian() {
	// Lance Obsidian via "open -a Obsidian" (fonctionne sur macOS)
	cmd := exec.Command("open", "-a", "Obsidian")
	if err := cmd.Start(); err != nil {
		fmt.Println("Erreur lors de l'ouverture d'Obsidian :", err)
		return
	}
	fmt.Println("Obsidian est lancé !")
}

func createNote(noteName string) {
	notePath := filepath.Join(vaultPath, noteName+".md")
	// Vérifier si la note existe déjà
	if _, err := os.Stat(notePath); err == nil {
		fmt.Println("La note existe déjà :", notePath)
		return
	} else if !os.IsNotExist(err) {
		fmt.Println("Erreur lors de la vérification de l'existence de la note :", err)
		return
	}
	file, err := os.Create(notePath)
	if err != nil {
		fmt.Println("Erreur lors de la création de la note :", err)
		return
	}
	defer file.Close()
	fmt.Println("Note créée :", notePath)
}

func writeNote(noteName string) {
	notePath := filepath.Join(vaultPath, noteName+".md")
	// Vérifier si la note existe
	if _, err := os.Stat(notePath); os.IsNotExist(err) {
		fmt.Println("La note n'existe pas :", notePath)
		return
	}
	// Ouvre nano pour éditer la note
	cmd := exec.Command("nano", notePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("Erreur lors de l'édition de la note :", err)
		return
	}
}

func deleteNote(noteName string) {
	notePath := filepath.Join(vaultPath, noteName+".md")
	if _, err := os.Stat(notePath); os.IsNotExist(err) {
		fmt.Println("La note n'existe pas :", notePath)
		return
	}
	if err := os.Remove(notePath); err != nil {
		fmt.Println("Erreur lors de la suppression de la note :", err)
		return
	}
	fmt.Println("Note supprimée :", notePath)
}

func listNotes() {
	entries, err := os.ReadDir(vaultPath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du coffre :", err)
		return
	}
	fmt.Println("Liste des notes :")
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			fmt.Println(" -", entry.Name())
		}
	}
}

func searchNote(keyword string) {
	entries, err := os.ReadDir(vaultPath)
	if err != nil {
		fmt.Println("Erreur lors de la lecture du coffre :", err)
		return
	}
	fmt.Printf("Notes contenant \"%s\" dans leur nom :\n", keyword)
	found := false
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			if strings.Contains(strings.ToLower(entry.Name()), strings.ToLower(keyword)) {
				fmt.Println(" -", entry.Name())
				found = true
			}
		}
	}
	if !found {
		fmt.Println("Aucune note ne correspond à la recherche.")
	}
}

// Met à jour le chemin du coffre et le sauvegarde dans le fichier de configuration
func updateVaultConfig(newPath string) {
	// Si le chemin n'existe pas, le créer
	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		if err := os.MkdirAll(newPath, 0755); err != nil {
			fmt.Println("Erreur lors de la création du nouveau coffre :", err)
			return
		}
		fmt.Println("Nouveau coffre créé dont le chemin est :", newPath)
	}
	vaultPath = newPath
	config := Config{VaultPath: newPath}
	if err := saveConfig(config); err != nil {
		fmt.Println("Erreur lors de la sauvegarde de la configuration :", err)
		return
	}
	fmt.Println("Le chemin du coffre a été mis à jour :", newPath)
}

func pushChanges() {
	// Ajouter tous les changements
	cmd := exec.Command("git", "-C", vaultPath, "add", ".")
	if err := cmd.Run(); err != nil {
		fmt.Println("Erreur lors de l'ajout des modifications :", err)
		return
	}
	// Effectuer le commit
	cmd = exec.Command("git", "-C", vaultPath, "commit", "-m", "Mise à jour via CLI")
	if err := cmd.Run(); err != nil {
		fmt.Println("Erreur lors du commit (Aucune modification ?) :", err)
	}
	// Pousser sur le dépôt distant
	cmd = exec.Command("git", "-C", vaultPath, "push")
	if err := cmd.Run(); err != nil {
		fmt.Println("Erreur lors du push :", err)
		return
	}
	fmt.Println("Modifications envoyées sur le dépôt Git avec succès !")
}
