# Gobsidian CLI

## Usage

## Exécution
Vous pouvez télécharger directement la dernière release et exécuter le CLI via terminal
```
./gobsidian
```

## Compilation
```
go build -o gobsidian main.go
./gobsidian
```

## Commandes disponibles
-  open              : Lancer Obsidian
-  new <nom>         : Créer une nouvelle note (vérifie si elle existe déjà)
-  write <nom>       : Éditer une note existante
-  delete <nom>      : Supprimer une note existante
-  list              : Afficher toutes les notes
-  search <mot_clé>  : Rechercher des notes contenant le mot-clé dans leur nom
-  vault             : Afficher l'emplacement actuel du coffre
-  setvault <chemin> : Changer et sauvegarder le chemin du coffre
-  push              : Envoyer les modifications en ligne (Git push)
-  help              : Afficher cette aide
-  quit ou exit      : Quitter le CLI
