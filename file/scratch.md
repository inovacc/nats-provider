​En Go, la estructura os.File representa un archivo abierto y proporciona métodos para interactuar con él. Aunque
os.File no es una interfaz, implementa varias interfaces estándar de Go, lo que permite su uso en una variedad de
contextos.​

Métodos de os.File
La estructura os.File ofrece una amplia gama de métodos para operaciones de lectura, escritura y manipulación de
archivos. Algunos de los métodos más comunes incluyen:​

Read(b []byte) (n int, err error)​
Google Groups
+5
GitHub
+5
Alex Edwards
+5

Write(b []byte) (n int, err error)​

Close() error​

Stat() (os.FileInfo, error)​
Go
+3
Googlesource
+3
Go Packages
+3

Seek(offset int64, whence int) (int64, error)​
GitHub

Sync() error​
Alex Edwards
+3
Googlesource
+3
GitHub
+3

Truncate(size int64) error​

Readdir(n int) ([]os.FileInfo, error)​

WriteString(s string) (n int, err error)​

Estos métodos permiten realizar operaciones comunes en archivos, como leer y escribir datos, mover el puntero de
lectura/escritura, sincronizar cambios con el almacenamiento y obtener información del archivo.​

Interfaces implementadas por os.File
Gracias a los métodos que implementa, os.File satisface varias interfaces estándar de Go, lo que permite su uso en
funciones que esperan estas interfaces:​

io.Reader: mediante el método Read​

io.Writer: mediante el método Write​

io.Closer: mediante el método Close​

io.Seeker: mediante el método Seek​

io.ReaderAt: mediante el método ReadAt​

io.WriterAt: mediante el método WriteAt​
GitHub

io.StringWriter: mediante el método WriteString​

Esto significa que puedes pasar un *os.File a funciones que aceptan estas interfaces, lo que proporciona una gran
flexibilidad en el manejo de archivos.