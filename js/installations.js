currentIndx=0;
MyImages=new Array();
MyImages[0]='img/sites/champak3kw.jpg';
MyImages[1]='img/sites/iisc12kw.jpg';
MyImages[2]='img/sites/ovh3kw.jpg';
MyImages[3]='img/sites/cmc40kw.jpg';
MyImages[4]='img/sites/sb3kw.jpg';
MyImages[5]='img/sites/gr1kw.jpg';
MyImages[6]='img/sites/gkp500w.jpg';
MyImages[7]='img/sites/cisco1mw.jpg';

Messages=new Array();
Messages[0]='Champak 3KW';
Messages[1]='IISc 12KW';
Messages[2]='OVH 3KW';
Messages[3]='CMC 40KW';
Messages[4]='SB 3KW';
Messages[5]='585 1KW';
Messages[6]='GKP 0.5KW';
Messages[7]='Cisco 1MW';

imagesPreloaded = new Array(8);

for (var i = 0; i < MyImages.length ; i++){
  imagesPreloaded[i] = new Image(640,480);
  imagesPreloaded[i].src=MyImages[i];
}
function writeImageNumber() {
  oSpan=document.getElementById("sp1");
  oSpan.innerHTML="Image "+eval(currentIndx+1)+" of "+MyImages.length;
}
function Nexter(){
  if (currentIndx<=imagesPreloaded.length-1){
    currentIndx=currentIndx+1;
    document.theImage.src=imagesPreloaded[currentIndx].src
    document.getElementById('text1').innerHTML=Messages[currentIndx];
  }
  else {
    currentIndx=0
    document.theImage.src=imagesPreloaded[currentIndx].src
    document.getElementById('text1').innerHTML=Messages[currentIndx];
    //document.form1.text1.value=Messages[currentIndx];
  }
  writeImageNumber();
}
function Backer(){
  if (currentIndx>0){
    currentIndx=currentIndx-1;
    document.theImage.src=imagesPreloaded[currentIndx].src
    document.getElementById('text1').innerHTML=Messages[currentIndx];
    //document.form1.text1.value=Messages[currentIndx];
  }
  else {
    currentIndx=7
    document.theImage.src=imagesPreloaded[currentIndx].src
    document.getElementById('text1').innerHTML=Messages[currentIndx];
  }
  writeImageNumber();
}

function automatically() {
  if (document.form1.automatic.checked) {
    if (currentIndx<imagesPreloaded.length){
      currentIndx=currentIndx
    }
    else {
      currentIndx=0
    }
    writeImageNumber()
    document.theImage.src=imagesPreloaded[currentIndx].src
    document.form1.text1.value=Messages[currentIndx];
    currentIndx=currentIndx+1;
    var delay = setTimeout("automatically()",3500)
  }
}
