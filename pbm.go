package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PBM représente la structure d'une image PBM
type PBM struct {
	data          [][]bool
	width, height int
	magicNumber   string
}

// ReadPBM lit une image PBM à partir d'un fichier et renvoie une structure qui représente l'image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var magicNumber string
	var width, height int
	var data [][]bool

	lineCount := 0
	for scanner.Scan() {
		line := strings.Fields(scanner.Text()) // Utiliser Fields pour diviser la ligne en mots
		if len(line) == 0 {
			continue
		}

		if lineCount == 0 {
			magicNumber = line[0]
		} else if lineCount == 1 {
			_, err := fmt.Sscanf(line[0], "%d", &width)
			if err != nil {
				return nil, fmt.Errorf("Erreur de format des dimensions: %v", err)
			}

			_, err = fmt.Sscanf(line[1], "%d", &height)
			if err != nil {
				return nil, fmt.Errorf("Erreur de format des dimensions: %v", err)
			}

			data = make([][]bool, height)
			for i := range data {
				data[i] = make([]bool, width)
			}
		} else {
			// Traitement des valeurs des pixels
			for j, val := range line {
				if val == "1" {
					data[lineCount-2][j] = true
				} else if val == "0" {
					data[lineCount-2][j] = false
				} else {
					return nil, fmt.Errorf("Erreur de format des valeurs des pixels: valeur inattendue %s", val)
				}
			}
		}
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &PBM{
		data:        data,
		width:       width,
		height:      height,
		magicNumber: magicNumber,
	}, nil
}

// Display affiche l'image PBM.
func (pbm *PBM) Display() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			if pbm.data[i][j] {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// boolToInt convertit une valeur booléenne en entier (0 ou 1).
func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Size retourne la largeur et la hauteur de l'image.
func (pbm *PBM) Size() (int, int) {
	return pbm.width, pbm.height
}

// At retourne la valeur du pixel aux coordonnées (x, y).
func (pbm *PBM) At(x, y int) bool {
	return pbm.data[y][x]
}

// Set définit la valeur du pixel aux coordonnées (x, y).
func (pbm *PBM) Set(x, y int, value bool) {
	pbm.data[y][x] = value
}

// Save enregistre l'image PBM dans un fichier et renvoie une erreur s'il y a un problème.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Écriture du magic number, de la largeur et de la hauteur
	fmt.Fprintf(writer, "%s\n%d %d\n", pbm.magicNumber, pbm.width, pbm.height)

	// Écriture des valeurs des pixels
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			if pbm.data[i][j] {
				fmt.Fprint(writer, "1")
			} else {
				fmt.Fprint(writer, "0")
			}
		}
		fmt.Fprintln(writer)
	}

	// Assurez-vous que toutes les données tamponnées sont écrites dans le fichier
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}

// Invert inverse les couleurs de l'image PBM.
func (pbm *PBM) Invert() {
	for i := 0; i < pbm.height; i++ {
		for j := 0; j < pbm.width; j++ {
			pbm.data[i][j] = !pbm.data[i][j]
		}
	}
}

// Flip inverse horizontalement l'image PBM.
func (pbm *PBM) Flip() {
	for i := 0; i < pbm.height; i++ {
		for j, k := 0, pbm.width-1; j < k; j, k = j+1, k-1 {
			pbm.data[i][j], pbm.data[i][k] = pbm.data[i][k], pbm.data[i][j]
		}
	}
}

// Flop inverse verticalement l'image PBM.
func (pbm *PBM) Flop() {
	for i, j := 0, pbm.height-1; i < j; i, j = i+1, j-1 {
		pbm.data[i], pbm.data[j] = pbm.data[j], pbm.data[i]
	}
}

// SetMagicNumber définit le magic number de l'image PBM.
func (pbm *PBM) SetMagicNumber(magicNumber string) {
	pbm.magicNumber = magicNumber
}

func main() {
	// Exemple d'utilisation
	image, err := ReadPBM("exemple.pbm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture de l'image PBM :", err)
		return
	}

	// Afficher les informations de l'image
	fmt.Printf("Nombre magique : %s\n", image.magicNumber)
	width, height := image.Size()
	fmt.Printf("Dimensions du dessin : %d x %d\n", width, height)

	// Afficher la valeur d'un pixel (par exemple, à la position (2, 3))
	value := image.At(2, 3)
	fmt.Printf("Valeur du pixel à la position (2, 3) : %t\n", value)

	// Appliquer des opérations (par exemple, Inverser, Flip, Flop)

	fmt.Println("Image normale:")
	image.Display()

	// Affiche l'image inversée
	image.Invert()
	fmt.Println("Image inversée:")
	image.Display()

	image.Flip()
	fmt.Println("Image inversée horizontalement:")
	image.Display()

	image.Flop()
	fmt.Println("Image inversée verticalement:")
	image.Display()
	// Enregistrer l'image modifiée
	err = image.Save("image_modifiee.pbm")
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement de l'image PBM :", err)
		return
	}
}