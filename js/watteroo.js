$(document).ready(function () {
	$('a.blank').click(function() {
        window.open(this.href);
        return false;
    });

});

function numberWithCommas(x) {
  return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g,",");
} 
	  
function get_started(){
      $('#dummy').remove();
      $('#dummy1').empty();
      $('#dummy1').load('whatresi.html');
    }

function view_products(){
      $('#dummy').remove();
      $('#dummy1').empty();
      $('#dummy1').load('catalog.html');
    }
function learn_more(){
      $('#dummy').remove();
      $('#dummy1').empty();
      $('#dummy1').load('whatresi.html');
    }

function rd(){
      $('#dummy').remove();
      $('#dummy1').empty();
      $('#dummy1').load('rd.html');
    }

function profile(){
      $('#dummy').remove();
      $('#dummy1').empty();
      $('#dummy1').load('company.html');
    }

function careers(){
      $('#dummy').remove();
      $('#dummy1').empty();
      $('#dummy1').load('careers.html');
    }

function login(){
      $('#dummy').remove();
      $('#dummy1').empty();
      $('#dummy1').load('login.html');
    }
function pri_pol(){
      $('#dummy').remove();
      $('#dummy1').empty();
      $('#dummy1').load('pri_pol.html');
    }
function terms(){
      $('#dummy').remove();
      $('#dummy1').empty();
      $('#dummy1').load('terms.html');
    }
