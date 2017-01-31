package main

import (
  "image"
  "image/png"
  "io/ioutil"
  "fmt"
  "path/filepath"
  "os"
  "strings"
)

func main(){
  image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
  css, err := os.OpenFile("icoshadow.css", os.O_WRONLY|os.O_CREATE, 0660);
  if err != nil {
    panic(err)
  }
  html, err := os.OpenFile("index.html", os.O_WRONLY|os.O_CREATE, 0660);
  if err != nil {
    panic(err)
  }

  files, _ := ioutil.ReadDir("./icons")
  if _, err = html.WriteString("<link rel='stylesheet' type='text/css' href='icoshadow.css'>"); err != nil {
    panic(err)
  }
  for _, icon := range files {
    if _, err = css.WriteString( processFile("icons/" + fmt.Sprintf("%v",icon.Name()))); err != nil {
      panic(err)
    }
    if _, err = html.WriteString(makeHTML(klass(fmt.Sprintf("%v",icon.Name())))); err != nil {
      panic(err)
    }
  }
  css.Close()
  html.Close()
}
func makeHTML(classname string) string {
  out := ""
  out = out + "<pre>&lt;i class='"+classname+"'>&lt;/i></pre>\n<div class='"+classname+"'></div>\n"
  return out
}
func processFile(fname string) string{
  fmt.Println("./" + fname)
  file, err := os.Open("./" + fname)
  if err != nil {
    panic("Failed to open image")
  }
  config, _, err := image.DecodeConfig(file)
  if err != nil {
    panic(err)
  }
  width := config.Width
  height := config.Height
  file.Seek(0,0)
  img, _, err := image.Decode(file)
  var shadows []string
  for y := 0; y < height; y++ {
     for x := 0; x < width; x++ {
       r, g, b, a := img.At(x, y).RGBA()
       shadows = append(shadows, rgbToBoxString(x,y,int(r/257), int(g/257), int(b/257), int(a/257)))
     }
  }
  file.Close()
  return makeCss(shadows, fname, height, width)
}

func rgbToBoxString(x int, y int, r int, g int, b int, a int) string {
  return fmt.Sprintf("%vpx %vpx 0 rgba(%v,%v,%v, %v)", x, y, r, g, b, a)
}

func klass(path string) string {
  out := strings.TrimSuffix(path, filepath.Ext(path))
  out = strings.TrimPrefix(out, "icons/")
  out = strings.Replace(out, ".","-",-1)
  return out
}

func makeCss(shadows []string, fname string, height int, width int) string {
  out := "."+ klass(fname) +"{\n"
    out = out + "  height: " + fmt.Sprintf("%v",height) + "px;\n"
    out = out + "  width: " + fmt.Sprintf("%v",width) + "px;\n"
    out = out + "  position: relative;\n"
  out = out + "}\n"
  out = out + "."+ klass(fname) +":before{\n"
    out = out + "  content: '';\n"
    out = out + "  display: block;\n"
    out = out + "  position: absolute;\n"
    out = out + "  top: 0;\n"
    out = out + "  left: 0;\n"
    out = out + "  width: 1px;\n  height: 1px;\n"
    out = out + "  box-shadow: " + strings.Join(shadows,", ") + ";\n"
  out = out + "}\n"
  return out
}
