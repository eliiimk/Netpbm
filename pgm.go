package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PGM représente une image PGM.
type PGM struct {
	data          [][]uint8
	width, height int
	magicNumber   string
	max           int
}

// Display affiche le dessin de l'image PGM dans la console.
func (pgm *PGM) Display() {
	for _, row := range pgm.data {
		for _, value := range row {
			fmt.Printf("%3d ", value)
		}
		fmt.Println()
	}
}

// ReadPGM lit une image PGM à partir d'un fichier et renvoie une structure qui représente l'image.
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var magicNumber string
	var width, height, max int

	// Lire les informations d'en-tête PGM.
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// Ignorer les commentaires
			continue
		}
		if magicNumber == "" {
			magicNumber = line
			if magicNumber != "P2" {
				return nil, fmt.Errorf("format PGM non pris en charge: %s", magicNumber)
			}
		} else {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				width, _ = strconv.Atoi(fields[0])
				height, _ = strconv.Atoi(fields[1])
				break // Sortir après avoir lu les dimensions
			}
		}
	}

	if width == 0 || height == 0 {
		return nil, fmt.Errorf("dimensions de l'image non spécifiées")
	}

	// Lire la valeur maximale
	scanner.Scan()
	max, _ = strconv.Atoi(scanner.Text())

	data := make([][]uint8, height)
	for i := 0; i < height; i++ {
		data[i] = make([]uint8, width)
	}

	// Lire les données de l'image.
	for i := 0; i < height; i++ {
		scanner.Scan()
		line := scanner.Text()
		if len(line) == 0 {
			return nil, fmt.Errorf("ligne vide à la position %d", i+2) // +2 pour compenser le 0-indexing et le fait que la première ligne est le chiffre magique "P2"
		}

		values := strings.Fields(line)

		if len(values) < width {
			return nil, fmt.Errorf("nombre insuffisant de valeurs sur la ligne %d, valeurs trouvées: %v", i+2, values) // +2 pour compenser le 0-indexing et le fait que la première ligne est le chiffre magique "P2"
		}

		// S'assurer que le tableau a une taille suffisante
		if len(data[i]) < width {
			data[i] = make([]uint8, width)
		}

		for j := 0; j < width; j++ {
			value, err := strconv.Atoi(values[j])
			if err != nil {
				return nil, err
			}
			data[i][j] = uint8(value)
		}
	}

	return &PGM{data, width, height, magicNumber, max}, nil
}

// Size renvoie la largeur et la hauteur de l'image.
func (pgm *PGM) Size() (int, int) {
	return pgm.width, pgm.height
}

// At renvoie la valeur du pixel en (x, y).
func (pgm *PGM) At(x, y int) uint8 {
	return pgm.data[y][x]
}

// Set définit la valeur du pixel à (x, y).
func (pgm *PGM) Set(x, y int, value uint8) {
	pgm.data[y][x] = value
}

// Save enregistre l'image PGM dans un fichier et renvoie une erreur en cas de problème.
func (pgm *PGM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "%s\n", pgm.magicNumber)
	fmt.Fprintf(writer, "%d %d\n", pgm.width, pgm.height)
	fmt.Fprintf(writer, "%d\n", pgm.max)

	for _, row := range pgm.data {
		for _, value := range row {
			fmt.Fprintf(writer, "%d ", value)
		}
		fmt.Fprintln(writer)
	}

	return nil
}

// Inverser inverse les couleurs de l'image PGM.
func (pgm *PGM) Invert() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j] = uint8(pgm.max) - pgm.data[i][j]
		}
	}
}

// Flip retourne l'image PGM horizontalement.
func (pgm *PGM) Flip() {
	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width/2; j++ {
			pgm.data[i][j], pgm.data[i][pgm.width-j-1] = pgm.data[i][pgm.width-j-1], pgm.data[i][j]
		}
	}
}

// Flop fait basculer l'image PGM verticalement.
func (pgm *PGM) Flop() {
	for i := 0; i < pgm.height/2; i++ {
		for j := 0; j < pgm.width; j++ {
			pgm.data[i][j], pgm.data[pgm.height-i-1][j] = pgm.data[pgm.height-i-1][j], pgm.data[i][j]
		}
	}
}

// SetMagicNumber définit le nombre magique de l'image PGM.
func (pgm *PGM) SetMagicNumber(magicNumber string) {
	pgm.magicNumber = magicNumber
}

// SetMaxValue définit la valeur maximale de l'image PGM.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.max = int(maxValue)
}

// Rotate90CW fait pivoter l'image PGM de 90° dans le sens des aiguilles d'une montre.
func (pgm *PGM) Rotate90CW() {
	rotatedData := make([][]uint8, pgm.width)
	for i := 0; i < pgm.width; i++ {
		rotatedData[i] = make([]uint8, pgm.height)
	}

	for i := 0; i < pgm.height; i++ {
		for j := 0; j < pgm.width; j++ {
			rotatedData[j][pgm.height-i-1] = pgm.data[i][j]
		}
	}

	pgm.data = rotatedData
	pgm.width, pgm.height = pgm.height, pgm.width
}

// ToPBM convertit l'image PGM en PBM.
func (pgm *PGM) ToPBM() *PBM {
	pbmData := make([][]bool, pgm.height)
	for i := 0; i < pgm.height; i++ {
		pbmData[i] = make([]bool, pgm.width)
		for j := 0; j < pgm.width; j++ {
			pbmData[i][j] = pgm.data[i][j] > uint8(pgm.max/2)
		}
	}

	return &PBM{pbmData, pgm.width, pgm.height}
}

// PBM représente une image PBM.
type PBM struct {
	data          [][]bool
	width, height int
}

// Save enregistre l'image PBM dans un fichier et renvoie une erreur en cas de problème.
func (pbm *PBM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "P1\n")
	fmt.Fprintf(writer, "%d %d\n", pbm.width, pbm.height)

	for _, row := range pbm.data {
		for _, value := range row {
			if value {
				fmt.Fprintf(writer, "1 ")
			} else {
				fmt.Fprintf(writer, "0 ")
			}
		}
		fmt.Fprintln(writer)
	}

	return nil
}

func main() {
	// Exemple d'utilisation
	pgm, err := ReadPGM("exemple.pgm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture de l'image PGM:", err)
		return
	}
	nouveauNombreMagique := "P2"
	pgm.SetMagicNumber(nouveauNombreMagique)
	fmt.Println()
	fmt.Println("Nouveau nombre magique:", pgm.magicNumber)

	width, height := pgm.Size()
	fmt.Println()
	fmt.Println("Taille de l'image:", width, height)

	// Afficher la nouvelle valeur maximale
	ValeurMax := uint8(10)
	pgm.SetMaxValue(ValeurMax)
	fmt.Println()
	fmt.Println("La valeur maximale de l'image:", pgm.max)

	x, y := 0, 9
	pixelValue := pgm.At(x, y)
	fmt.Println()
	fmt.Printf("Valeur du pixel à la position (%d, %d): %d\n", x, y, pixelValue)

	newPixelValue := uint8(4) // Nouvelle valeur que vous souhaitez définir
	pgm.Set(x, y, newPixelValue)
	fmt.Println()
	fmt.Printf("Nouvelle valeur du pixel à la position (%d, %d): %d\n", x, y, pgm.At(x, y))

	fmt.Println()
	fmt.Println("Image normale :")
	fmt.Println()
	pgm.Display()
	fmt.Println()

	pgm.Invert()
	fmt.Println("Image avec les couleurs inversé :")
	fmt.Println()
	pgm.Display()
	fmt.Println()

	pgm.Flip()
	fmt.Println("Image retourner à l'horizontalement :")
	fmt.Println()
	pgm.Display()
	fmt.Println()

	pgm.Flop()
	fmt.Println("Image renverser à la verticalement :")
	fmt.Println()
	pgm.Display()
	fmt.Println()

	pgm.Rotate90CW()
	fmt.Println("Image tourné à 90° dans le sens des aiguilles d'une montre :")
	fmt.Println()
	pgm.Display()
	fmt.Println()

	// Sauvegarde de l'image modifiée
	err = pgm.Save("image_modifiee.pgm")
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement de l'image PGM modifiée:", err)
		return
	}

	// Conversion en PBM
	pbm := pgm.ToPBM()

	// Affichage de l'image PBM
	pbmWidth, pbmHeight := pbm.width, pbm.height
	fmt.Println()
	fmt.Println("Taille de l'image PBM:", pbmWidth, pbmHeight)
	fmt.Println()
	fmt.Println("Image PBM:")
	fmt.Println()
	for _, row := range pbm.data {
		for _, value := range row {
			if value {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}

	// Sauvegarde de l'image PBM
	err = pbm.Save("image_pbm.pbm")
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement de l'image PBM:", err)
		return
	}
}