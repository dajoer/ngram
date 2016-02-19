# ngram

Ngram ist ein N-Gramm-Sprachmodell. 

## Verwendung

Zum lernen muss das Programm mit dem Flag --learn und einem Dateinamen als Argument ausgeführt werden. Das Programm nicht auf Stdin Zeilenweise Sätze entgegen und speichert das Gelernte in der angegebenen Datei.

Wird das Programm nur mit einem Dateinamen ausgeführt, nimmt es Zeilenweise Sätze entgegen und gibt die Satzwahrscheinlichkeit zurück. Ist der Flag -b gesetzt werden bis zum EOF Signal Sätze eingelesen und nur der Wahrscheinlichste zurückgegeben. Der Flag -v sorgt dafür, dass sowohl der Satz, als auch die Satzwahrscheinlichkeit zurückgegeben wird.

Beispiele zum lernen und für eine einfache Anwendung sind in den Bash-Scripts learn.sh und example_verbose.sh gegeben. Die Scripts benötigen als Argument die Datei, in der die Daten des Sprachmodells gespeichert sind. learn.sh erstellt die Datei, falls sie noch nicht vorhanden ist. Das Script example_verbose.sh benötigt das Programm heapsAlg (https://github.com/dajoer/heapsAlg)
