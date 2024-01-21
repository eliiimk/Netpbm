package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PPM représente une image PPM.
type PPM struct {
	data          [][][]uint8
	width, height int
	magicNumber   string
	max           int
}

type Pixel struct {
	Red, Green, Blue uint8
}

// Display affiche le dessin de l'image PPM dans la console.
func (ppm *PPM) Display() {
	for _, row := range ppm.data {
		for _, pixel := range row {
			fmt.Printf("%3d %3d %3d ", pixel[0], pixel[1], pixel[2])
		}
		fmt.Println()
	}
}

// ReadPPM lit une image PPM à partir d'un fichier et renvoie une structure qui représente l'image.
func ReadPPM(filename string) (*PPM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var magicNumber string
	var width, height, max int

	// Lire les informations d'en-tête PPM.
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			// Ignorer les commentaires
			continue
		}
		if magicNumber == "" {
			magicNumber = line
			if magicNumber != "P3" {
				return nil, fmt.Errorf("format PPM non pris en charge: %s", magicNumber)
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

	data := make([][][]uint8, height)
	for i := 0; i < height; i++ {
		data[i] = make([][]uint8, width)
	}

	// Lire les données de l'image.
	for i := 0; i < height; i++ {
		scanner.Scan()
		line := scanner.Text()
		if len(line) == 0 {
			return nil, fmt.Errorf("ligne vide à la position %d", i+2) // +2 pour compenser le 0-indexing et le fait que la première ligne est le chiffre magique "P3"
		}

		values := strings.Fields(line)

		if len(values) < width*3 {
			return nil, fmt.Errorf("nombre insuffisant de valeurs sur la ligne %d, valeurs trouvées: %v", i+2, values) // +2 pour compenser le 0-indexing et le fait que la première ligne est le chiffre magique "P3"
		}

		// S'assurer que le tableau a une taille suffisante
		if len(data[i]) < width {
			data[i] = make([][]uint8, width)
		}

		for j := 0; j < width; j++ {
			r, _ := strconv.Atoi(values[j*3])
			g, _ := strconv.Atoi(values[j*3+1])
			b, _ := strconv.Atoi(values[j*3+2])
			data[i][j] = []uint8{uint8(r), uint8(g), uint8(b)}
		}
	}

	return &PPM{data, width, height, magicNumber, max}, nil
}

// Size renvoie la largeur et la hauteur de l'image.
func (ppm *PPM) Size() (int, int) {
	return ppm.width, ppm.height
}

// At renvoie la valeur du pixel en (x, y).
func (ppm *PPM) At(x, y int) []uint8 {
	return ppm.data[y][x]
}

// / Save enregistre l'image PPM dans un fichier et renvoie une erreur en cas de problème.
func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	fmt.Fprintf(writer, "%s\n", ppm.magicNumber)
	fmt.Fprintf(writer, "%d %d\n", ppm.width, ppm.height)
	fmt.Fprintf(writer, "%d\n", ppm.max)

	for _, row := range ppm.data {
		for _, pixel := range row {
			fmt.Fprintf(writer, "%d %d %d ", pixel[0], pixel[1], pixel[2])
		}
		fmt.Fprintln(writer)
	}

	return nil
}

// Inverser inverse les couleurs de l'image PPM.
func (ppm *PPM) Invert() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			for k := 0; k < 3; k++ {
				ppm.data[i][j][k] = uint8(ppm.max) - ppm.data[i][j][k]
			}
		}
	}
}

// Flip retourne l'image PPM horizontalement.
func (ppm *PPM) Flip() {
	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width/2; j++ {
			ppm.data[i][j], ppm.data[i][ppm.width-j-1] = ppm.data[i][ppm.width-j-1], ppm.data[i][j]
		}
	}
}

// Flop fait basculer l'image PPM verticalement.
func (ppm *PPM) Flop() {
	for i := 0; i < ppm.height/2; i++ {
		for j := 0; j < ppm.width; j++ {
			ppm.data[i][j], ppm.data[ppm.height-i-1][j] = ppm.data[ppm.height-i-1][j], ppm.data[i][j]
		}
	}
}

// SetMagicNumber définit le nombre magique de l'image PPM.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.magicNumber = magicNumber
}

// SetMaxValue définit la valeur maximale de l'image PPM.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	ppm.max = int(maxValue)
}

// Rotate90CW fait pivoter l'image PPM de 90° dans le sens des aiguilles d'une montre.
func (ppm *PPM) Rotate90CW() {
	rotatedData := make([][][]uint8, ppm.width)
	for i := 0; i < ppm.width; i++ {
		rotatedData[i] = make([][]uint8, ppm.height)
	}

	for i := 0; i < ppm.height; i++ {
		for j := 0; j < ppm.width; j++ {
			rotatedData[j][ppm.height-i-1] = ppm.data[i][j]
		}
	}

	ppm.data = rotatedData
	ppm.width, ppm.height = ppm.height, ppm.width
}

// Point représente un point dans l'image.
type Point struct {
	X, Y int
}

// Fonction utilitaire abs pour obtenir la valeur absolue d'un nombre entier.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// DrawLine trace une ligne entre deux points.
func (ppm *PPM) DrawLine(p1, p2 Point, couleur Pixel) {
	x1, y1 := p1.X, p1.Y
	x2, y2 := p2.X, p2.Y

	dx := x2 - x1
	dy := y2 - y1

	var steps int
	if abs(dx) > abs(dy) {
		steps = abs(dx)
	} else {
		steps = abs(dy)
	}

	xInc := float64(dx) / float64(steps)
	yInc := float64(dy) / float64(steps)

	x, y := float64(x1), float64(y1)

	for i := 0; i <= steps; i++ {
		// Convertir les coordonnées en entiers et vérifier les limites
		xInt, yInt := int(x), int(y)
		if xInt >= 0 && xInt < ppm.width && yInt >= 0 && yInt < ppm.height {
			ppm.Set(xInt, yInt, []uint8{couleur.Red, couleur.Green, couleur.Blue})
		}
		x += xInc
		y += yInc
	}
}

// ...

// Set définit la valeur du pixel à (x, y).
func (ppm *PPM) Set(x, y int, value []uint8) {
	// Assurez-vous que ppm.data[y] a une longueur suffisante
	for len(ppm.data) <= y {
		ppm.data = append(ppm.data, make([][]uint8, ppm.width))
	}

	// Assurez-vous que ppm.data[y][x] a une longueur suffisante
	for len(ppm.data[y]) <= x {
		ppm.data[y] = append(ppm.data[y], make([]uint8, 3))
	}

	ppm.data[y][x] = value
}

// DrawTriangle dessine un triangle.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, couleur Pixel) {
	ppm.DrawLine(p1, p2, couleur)
	ppm.DrawLine(p2, p3, couleur)
	ppm.DrawLine(p3, p1, couleur)
}

// drawLine dessine une ligne entre deux points.
func (ppm *PPM) drawLine(start, end Point, color Pixel) {
	x0, y0 := start.X, start.Y
	x1, y1 := end.X, end.Y

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

	var sx, sy int

	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}

	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		ppm.setPixel(x0, y0, color)

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err

		if e2 > -dy {
			err -= dy
			x0 += sx
		}

		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// DrawFilledTriangle dessine un triangle rempli dans l'image PPM.
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	// Utiliser l'algorithme de tracé de ligne pour dessiner les trois côtés du triangle.
	ppm.drawFilledLine(p1, p2, color)
	ppm.drawFilledLine(p2, p3, color)
	ppm.drawFilledLine(p3, p1, color)
}

// drawFilledLine dessine une ligne remplie entre deux points.
func (ppm *PPM) drawFilledLine(start, end Point, color Pixel) {
	x0, y0 := start.X, start.Y
	x1, y1 := end.X, end.Y

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

	var sx, sy int

	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}

	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		// Dessiner la ligne horizontale entre les deux points.
		ppm.drawHorizontalLine(y0, x0, x1, color)

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err

		if e2 > -dy {
			err -= dy
			x0 += sx
		}

		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// DrawPolygon dessine un polygone dans l'image PPM.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	// Vérifier que la liste de points n'est pas vide.
	if len(points) < 3 {
		fmt.Println("Un polygone doit avoir au moins trois points.")
		return
	}

	// Utiliser l'algorithme de tracé de ligne pour dessiner les côtés du polygone.
	for i := 0; i < len(points)-1; i++ {
		ppm.drawLine(points[i], points[i+1], color)
	}
	// Dessiner la dernière ligne reliant le dernier point au premier point.
	ppm.drawLine(points[len(points)-1], points[0], color)
}

// DrawFilledPolygon dessine un polygone rempli dans l'image PPM.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// Vérifier que la liste de points n'est pas vide.
	if len(points) < 3 {
		fmt.Println("Un polygone rempli doit avoir au moins trois points.")
		return
	}

	// Trouver la boîte englobante du polygone pour délimiter la zone à remplir.
	minX, minY, maxX, maxY := findBoundingBox(points)

	// Itérer sur chaque ligne de la boîte englobante et remplir les pixels à l'intérieur du polygone.
	for y := minY; y <= maxY; y++ {
		intersections := findIntersections(points, y)

		// Dessiner une ligne horizontale entre chaque paire d'intersections.
		for i := 0; i < len(intersections)-1; i += 2 {
			startX := max(intersections[i], minX)
			endX := min(intersections[i+1], maxX)

			ppm.drawHorizontalLine(y, startX, endX, color)
		}
	}
}

// drawHorizontalLine dessine une ligne horizontale entre les coordonnées spécifiées.
func (ppm *PPM) drawHorizontalLine(y, startX, endX int, color Pixel) {
	for x := startX; x <= endX; x++ {
		ppm.setPixel(x, y, color)
	}
}

// DrawFilledRectangle dessine un rectangle rempli dans l'image PPM.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	// Vérifier que les coordonnées du point ne dépassent pas les dimensions de l'image.
	if p1.X < 0 || p1.X >= ppm.width || p1.Y < 0 || p1.Y >= ppm.height {
		fmt.Println("Les coordonnées du point sont hors des limites de l'image.")
		return
	}

	// Assurer que la largeur et la hauteur du rectangle sont positives.
	if width <= 0 || height <= 0 {
		fmt.Println("La largeur et la hauteur du rectangle doivent être positives.")
		return
	}

	// Vérifier que le rectangle ne dépasse pas les limites de l'image.
	if p1.X+width > ppm.width || p1.Y+height > ppm.height {
		fmt.Println("Le rectangle dépasse les limites de l'image.")
		return
	}

	// Dessiner le rectangle rempli.
	for i := p1.Y; i < p1.Y+height; i++ {
		for j := p1.X; j < p1.X+width; j++ {
			ppm.data[i][j] = []uint8{color.Red, color.Green, color.Blue}
		}
	}
}

// DrawCircle dessine un cercle dans l'image PPM.
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	// Vérifier que les coordonnées du centre ne dépassent pas les dimensions de l'image.
	if center.X < 0 || center.X >= ppm.width || center.Y < 0 || center.Y >= ppm.height {
		fmt.Println("Les coordonnées du centre sont hors des limites de l'image.")
		return
	}

	// Assurer que le rayon du cercle est positif.
	if radius <= 0 {
		fmt.Println("Le rayon du cercle doit être positif.")
		return
	}

	// Vérifier que le cercle ne dépasse pas les limites de l'image.
	if center.X-radius < 0 || center.X+radius >= ppm.width || center.Y-radius < 0 || center.Y+radius >= ppm.height {
		fmt.Println("Le cercle dépasse les limites de l'image.")
		return
	}

	// Dessiner le cercle en utilisant l'algorithme de tracé de cercle de Bresenham.
	x := radius
	y := 0
	decision := 1 - x

	for y <= x {
		ppm.setCirclePixels(center, x, y, color)

		y++
		if decision <= 0 {
			decision += 2*y + 1
		} else {
			x--
			decision += 2*(y-x) + 1
		}
	}
}

// setCirclePixels définit les pixels du cercle symétriquement autour du centre.
func (ppm *PPM) setCirclePixels(center Point, x, y int, color Pixel) {
	ppm.setPixel(center.X+x, center.Y+y, color)
	ppm.setPixel(center.X-x, center.Y+y, color)
	ppm.setPixel(center.X+x, center.Y-y, color)
	ppm.setPixel(center.X-x, center.Y-y, color)
	ppm.setPixel(center.X+y, center.Y+x, color)
	ppm.setPixel(center.X-y, center.Y+x, color)
	ppm.setPixel(center.X+y, center.Y-x, color)
	ppm.setPixel(center.X-y, center.Y-x, color)
}

// setPixel définit la couleur d'un pixel dans le cercle.
func (ppm *PPM) setPixel(x, y int, color Pixel) {
	// Assurez-vous que les coordonnées sont dans les limites de l'image.
	if x >= 0 && x < ppm.width && y >= 0 && y < ppm.height {
		ppm.data[y][x] = []uint8{color.Red, color.Green, color.Blue}
	}
}

// DrawFilledCircle dessine un cercle rempli dans l'image PPM.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	// Vérifier que les coordonnées du centre ne dépassent pas les dimensions de l'image.
	if center.X < 0 || center.X >= ppm.width || center.Y < 0 || center.Y >= ppm.height {
		fmt.Println("Les coordonnées du centre sont hors des limites de l'image.")
		return
	}

	// Assurer que le rayon du cercle est positif.
	if radius <= 0 {
		fmt.Println("Le rayon du cercle doit être positif.")
		return
	}

	// Vérifier que le cercle ne dépasse pas les limites de l'image.
	if center.X-radius < 0 || center.X+radius >= ppm.width || center.Y-radius < 0 || center.Y+radius >= ppm.height {
		fmt.Println("Le cercle dépasse les limites de l'image.")
		return
	}

	// Dessiner le cercle rempli en utilisant l'algorithme de tracé de cercle de Bresenham.
	x := radius
	y := 0
	decision := 1 - x

	for y <= x {
		ppm.drawHorizontalLine(center.Y+y, center.X-x, center.X+x, color)
		ppm.drawHorizontalLine(center.Y-y, center.X-x, center.X+x, color)

		y++
		if decision <= 0 {
			decision += 2*y + 1
		} else {
			x--
			decision += 2*(y-x) + 1
		}
	}
}

// findBoundingBox trouve la boîte englobante d'un polygone.
func findBoundingBox(points []Point) (minX, minY, maxX, maxY int) {
	// Initialiser avec les premières coordonnées.
	minX, minY, maxX, maxY = points[0].X, points[0].Y, points[0].X, points[0].Y

	// Mettre à jour les coordonnées min/max en parcourant tous les points.
	for _, point := range points {
		if point.X < minX {
			minX = point.X
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	return minX, minY, maxX, maxY
}

// findIntersections trouve les intersections d'une ligne horizontale avec les côtés du polygone.
func findIntersections(points []Point, y int) []int {
	intersections := make([]int, 0)

	// Itérer sur chaque côté du polygone.
	for i := 0; i < len(points); i++ {
		j := (i + 1) % len(points) // Indice du prochain point (ferme la boucle)

		// Vérifier si la ligne horizontale intersecte le côté du polygone.
		if (points[i].Y > y && points[j].Y <= y) || (points[j].Y > y && points[i].Y <= y) {
			// Calculer l'intersection en utilisant l'équation de la droite.
			x := int(float64(points[i].X) + float64(y-points[i].Y)/float64(points[j].Y-points[i].Y)*float64(points[j].X-points[i].X))

			// Ajouter l'intersection à la liste.
			intersections = append(intersections, x)
		}
	}

	return intersections
}

// min renvoie le minimum de deux entiers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max renvoie le maximum de deux entiers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Copy crée une copie de l'image PPM.
func (ppm *PPM) Copy() *PPM {
	copyData := make([][][]uint8, ppm.height)
	for i := range copyData {
		copyData[i] = make([][]uint8, ppm.width)
		for j := range copyData[i] {
			copyData[i][j] = make([]uint8, 3)
			copy(copyData[i][j], ppm.data[i][j])
		}
	}

	return &PPM{
		data:        copyData,
		width:       ppm.width,
		height:      ppm.height,
		magicNumber: ppm.magicNumber,
		max:         ppm.max,
	}
}

func main() {
	// Exemple d'utilisation
	ppm, err := ReadPPM("exemple.ppm")
	if err != nil {
		fmt.Println("Erreur lors de la lecture de l'image PPM:", err)
		return
	}
	nouveauNombreMagique := "P3"
	ppm.SetMagicNumber(nouveauNombreMagique)
	fmt.Println()
	fmt.Println("Nouveau nombre magique:", ppm.magicNumber)

	width, height := ppm.Size()
	fmt.Println()
	fmt.Println("Taille de l'image:", width, height)

	x, y := 0, 9
	pixelValue := ppm.At(x, y)
	fmt.Println()
	fmt.Printf("Valeur du pixel à la position (%d, %d): %v\n", x, y, pixelValue)

	newPixelValue := []uint8{100, 150, 200} // Nouvelle valeur que vous souhaitez définir
	ppm.Set(x, y, newPixelValue)
	fmt.Println()
	fmt.Printf("Nouvelle valeur du pixel à la position (%d, %d): %v\n", x, y, ppm.At(x, y))

	fmt.Println()
	fmt.Println("Image normale :")
	fmt.Println()
	ppm.Display()
	fmt.Println()

	ppm.Invert()
	fmt.Println("Image avec les couleurs inversées :")
	fmt.Println()
	ppm.Display()
	fmt.Println()

	ppm.Flip()
	fmt.Println("Image retournée horizontalement :")
	fmt.Println()
	ppm.Display()
	fmt.Println()

	ppm.Flop()
	fmt.Println("Image renversée verticalement :")
	fmt.Println()
	ppm.Display()
	fmt.Println()

	ppm.Rotate90CW()
	fmt.Println("Image tournée à 90° dans le sens des aiguilles d'une montre :")
	fmt.Println()
	ppm.Display()
	fmt.Println()

	// Création d'une nouvelle instance de Point avec de nouvelles valeurs
	newPoint := Point{X: 3, Y: 5}
	newPixel := []uint8{100, 150, 200}

	// Utilisation de la nouvelle instance de Point pour définir un nouveau pixel
	ppm.Set(newPoint.X, newPoint.Y, newPixel)
	fmt.Printf("Nouvelle valeur du pixel à la position (%d, %d): %v\n", newPoint.X, newPoint.Y, ppm.At(newPoint.X, newPoint.Y))
	fmt.Println()

	point1 := Point{X: 10, Y: 10}
	point2 := Point{X: 50, Y: 30}

	couleurLigne := Pixel{Red: 255, Green: 0, Blue: 0} // Rouge

	// Tracez une ligne entre les deux points
	ppm.DrawLine(point1, point2, couleurLigne)
	fmt.Println("ligne :")
	fmt.Println()
	ppm.Display()
	fmt.Println()

	// Couleur du rectangle
	couleurRectangle := Pixel{Red: 0, Green: 0, Blue: 255} // Bleu

	// Dessiner un rectangle rempli dans l'image PPM
	ppm.DrawFilledRectangle(point1, width, height, couleurRectangle)

	// Affichage de l'image avec le rectangle rempli
	fmt.Println("Rectangle rempli :")
	fmt.Println()
	ppm.Display()
	fmt.Println()

	// Dessiner un triangle
	p1 := Point{5, 5}
	p2 := Point{10, 15}
	p3 := Point{15, 5}
	ppm.DrawTriangle(p1, p2, p3, Pixel{255, 0, 0})

	// Affichage de l'image avec la ligne dessinée
	fmt.Println("triangle:")
	fmt.Println()
	ppm.Display()

	couleurTriangle := Pixel{Red: 0, Green: 255, Blue: 0} // Vert
	ppm.DrawFilledTriangle(p1, p2, p3, couleurTriangle)

	// Affichage de l'image avec le triangle rempli
	fmt.Println("Triangle rempli :")
	fmt.Println()
	ppm.Display()
	fmt.Println()

	centreCercle := Point{X: 30, Y: 30}
	rayonCercle := 15
	couleurCercle := Pixel{Red: 255, Green: 0, Blue: 0} // Rouge
	ppmCercle := ppm.Copy()                             // Créer une copie de l'image originale
	ppmCercle.DrawCircle(centreCercle, rayonCercle, couleurCercle)

	// Affichage de l'image avec le cercle
	fmt.Println("Image avec le cercle :")
	ppmCercle.Display()
	fmt.Println()

	// Dessiner un cercle rempli
	centreCercleRempli := Point{X: 80, Y: 30}
	rayonCercleRempli := 10
	couleurCercleRempli := Pixel{Red: 0, Green: 255, Blue: 0} // Vert
	ppmCercleRempli := ppm.Copy()                             // Créer une copie de l'image originale
	ppmCercleRempli.DrawFilledCircle(centreCercleRempli, rayonCercleRempli, couleurCercleRempli)

	// Affichage de l'image avec le cercle rempli
	fmt.Println("Image avec le cercle rempli :")
	ppmCercleRempli.Display()
	fmt.Println()

	// Dessiner un polygone
	pointsPolygone := []Point{
		{X: 10, Y: 60},
		{X: 20, Y: 80},
		{X: 30, Y: 60},
		{X: 25, Y: 50},
	}
	couleurPolygone := Pixel{Red: 0, Green: 0, Blue: 255} // Bleu
	ppmPolygone := ppm.Copy()                             // Créer une copie de l'image originale
	ppmPolygone.DrawPolygon(pointsPolygone, couleurPolygone)

	// Affichage de l'image avec le polygone
	fmt.Println("Image avec le polygone :")
	ppmPolygone.Display()
	fmt.Println()

	// Dessiner un polygone rempli
	pointsPolygoneRempli := []Point{
		{X: 60, Y: 60},
		{X: 70, Y: 80},
		{X: 80, Y: 60},
		{X: 75, Y: 50},
	}
	couleurPolygoneRempli := Pixel{Red: 255, Green: 255, Blue: 0} // Jaune
	ppmPolygoneRempli := ppm.Copy()                               // Créer une copie de l'image originale
	ppmPolygoneRempli.DrawFilledPolygon(pointsPolygoneRempli, couleurPolygoneRempli)

	// Affichage de l'image avec le polygone rempli
	fmt.Println("Image avec le polygone rempli :")
	ppmPolygoneRempli.Display()
	fmt.Println()

	// Sauvegarde de l'image modifiée
	err = ppm.Save("image_modifiee.ppm")
	if err != nil {
		fmt.Println("Erreur lors de l'enregistrement de l'image PPM modifiée:", err)
		return
	}
}