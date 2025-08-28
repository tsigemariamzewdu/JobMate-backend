package pdfparser

import (
	"fmt"
	"strings"
	"github.com/ledongthuc/pdf"
)

//parse pdf extracts text from  a pdf file 
//it returns the raw text with spaces adjusted between words

func ParsePDF(path string)(string,error){
	//open the pdf file
	f,r,err:=pdf.Open(path)

	if err !=nil{
		return "",fmt.Errorf("failed to open PDF: %v",err)
	}
	defer f.Close() //ensure file is closed when function exits

	var textBuilder strings.Builder
	
	//total number of pages in the pdf
	totalPage:=r.NumPage()

	//loop through each page
	for pageIndex:=1;pageIndex<=totalPage;pageIndex++{
		p:=r.Page(pageIndex)
		//skip empty pages
		if p.V.IsNull(){
			continue
		}

		//get text organized by rows
		rows,err:=p.GetTextByRow()
		if err!=nil{
			return "",fmt.Errorf("failed to extract text from page %d:%v",pageIndex,err)

		}
		if len(rows)==0{
			continue
		}

		//loop through each row
		for _,row:=range rows{
			var rowText strings.Builder
			prevX:=-1.0

			for _,word:=range row.Content{
				if prevX>=0 && word.X-prevX>1.5{
					rowText.WriteString(" ")
				}
				rowText.WriteString(word.S)  //append word
				prevX=word.X+word.W
			}

			//append the row to the page text trim the extra spaces

			textBuilder.WriteString(strings.TrimSpace(rowText.String())+"\n")
		}

		textBuilder.WriteString("\n")
	
	}
	return textBuilder.String(),nil
}