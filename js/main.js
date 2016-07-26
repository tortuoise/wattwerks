(function($){
	var tfToolTip = function(element, options){
		var settings = $.extend({}, $.fn.tfToolTip.defaults, options);
		var element = $(element);
		var template = settings.template.replace('{title}', (settings.title == '') ? element.attr('title') : settings.title);
		var tooltip = $(template);
		
		element.live('mouseenter', function(){
			
			var Offset = element.offset();
			var Height = element.outerHeight();
			var Width = element.outerWidth();
			
			tooltip.hide();
			tooltip.appendTo('body');
			tooltip.css({left: Offset.left + Width / 2 - tooltip.outerWidth() / 2, top: Offset.top - tooltip.outerHeight() - 10});
			tooltip.css({display: 'block', opacity: 0});
			tooltip.stop().animate({opacity: 1, top: Offset.top - tooltip.outerHeight()}, 200);
		});
		
		element.live('mouseleave', function(){
			tooltip.stop().fadeOut(200, function(){
				tooltip.remove();	
			});
		});
	};
	
	$.fn.tfToolTip = function(options){
        return this.each(function(key, value){
            var element = $(this);
			
            if (element.data('tftooltip')) return element.data('tftooltip');

            var tftooltip = new tfToolTip(this, options);

            element.data('tftooltip', tftooltip);
        });
	};
	
	$.fn.tfToolTip.defaults = {
		title: '',
		template: '<div class="tf-tooltip">{title}</div>'
	};

	$(document).ready(function(){
	  $('.tftooltip').tfToolTip();
	  $('#button-cart').bind('click', function() {
	    var product_id = $('#product_id').val();
	    var quantity = $('#selected_quantity').val();
	    $.ajax({
		    url: '/Cart',
		    type: 'post',
		    data: 'action=add' + '&product_id=' + product_id + '&quantity=' + quantity,
		    dataType: 'json',
		    success: function(json) {
		      if (json['redirect']) {
			location = json['redirect'];
		      }
		      
		      if (json['success']) {
			Notification(json['success'], 'success');
			
			var matches = /^(.*) - ([^ ]*)$/i.exec(json['total']);
			
			$('#cart .cart-count').text(matches[1]);
			$('#cart .cart-total-text').text(numberWithCommas(matches[2]));
		      }	
		    }
	    });
	  });
	  $('.update_cart').bind('click', function() {
	    var product_id = $(this).next().val();
	    var quantity = $(this).siblings().filter('#update_qntty').val();
	    $.ajax({
		    url: '/Cart',
		    type: 'post',
		    data: 'action=update' + '&product_id=' + product_id + '&quantity=' + quantity,
		    dataType: 'json',
		    success: function(json) {
		      if (json['redirect']) {
			location = json['redirect'];
		      }
		      
		      if (json['success']) {
                        window.location.replace("/Cart?action=checkout");
		      }	
		    }
	    });
	  });
	  $('.remove_cart').bind('click', function() {
	    var product_id = $(this).siblings().filter('#product_id').val();
	    $.ajax({
		    url: '/Cart',
		    type: 'post',
		    data: 'action=remove' + '&product_id=' + product_id,
		    dataType: 'json',
		    success: function(json) {
		      if (json['redirect']) {
			location = json['redirect'];
		      }
		      
		      if (json['success']) {
                        window.location.replace("/Cart?action=checkout");
		      }	
		    }
	    });
	  });
	  $('.checkout_cart').bind('click', function() {
	    $.ajax({
		    url: '/Cart',
		    type: 'post',
		    data: 'action=checkout',
		    dataType: 'json',
		    success: function(json) {
		      if (json['redirect']) {
			location = json['redirect'];
		      }
		      
		      if (json['success']) {
                        console.log(json['success']);
                        window.location.replace(json['success']);
		      }	
		    }
	    });
	  });
	  $('#button-review').bind('click', function() {
	    //var rating_radios = $("#ratings input:radio[name='rating']");
	    //var review_text = $("#review input:text[name='review']");
	    var product_id = $('#product_id').val();
	    var rating = $('[name=rating]:checked').val();//rating_radios.filter(':checked').val();
	    var review = $('textarea#review').val();
	    $.ajax({
		    url: '/products/review',
		    type: 'post',
		    data: '&product_id=' + product_id + '&rating=' + rating + '&review=' + review,
		    dataType: 'json',
		    success: function(json) {
		      if (json['failure']) {
			Notification(json['failure'], 'failure');
		      }
		      
		      if (json['success']) {
			Notification(json['success'], 'success');
		      }	
		    }
	    });
	  });
	});

})(jQuery);

function DeleteFromCart(product_id, quantity){
      /*if(getURLVar('route') == 'checkout/cart' || getURLVar('route') == 'checkout/checkout'){
		location = 'index.php?route=checkout/cart&remove=' + product_id;
      }else*/{

	quantity = typeof(quantity) != 'undefined' ? quantity : 1;
	$.ajax({
		url: '/Cart',
		type: 'post',
		data: 'action=remove' + '&product_id=' + product_id + '&quantity=' + quantity,
		dataType: 'json',
		success: function(json) {
		  if (json['redirect']) {
		    location = json['redirect'];
		  }
		  
		  if (json['success']) {
		    Notification(json['success'], 'success');
		    
		    var matches = /^(.*) - ([^ ]*)$/i.exec(json['total']);
		    
		    $('#cart .cart-count').text(matches[1]);
		    $('#cart .cart-total-text').text(numberWithCommas(matches[2]));
		  }	
		}
	});
	$('#cart .drop-body').load('/Cart' , function(){
		$('#cart').trigger('mouseenter').find('.cart-total').stop(true, true);
		$('#cart').find('.drop-content, .drop-arrow').show();
	});
      }
}

function Notification(msg, type){
	var box = $('<div class="' + type + '"><div class="close"></div><div class="icon"></div>' + msg + '</div>').hide();
	
	$('#notification').prepend(box);
	
	if(type != 'warning'){	
		box.fadeIn(500).delay(10000).fadeOut(500, function(){
			box.remove();	
		});
	}else{
		box.fadeIn(500);	
	}
	
	box.find('.close').on('click', function(){
		box.stop().fadeOut(500, function(){
			box.remove();	
		});	
	});
}

$(document).ready(function(){
  $(document).click(function(){
    $('.drop-content .drop-body:visible').slideUp(200, function(){
      $(this).parents('.drop-content').parent().removeClass('active').trigger('mouseleave').find('.drop-content, .drop-arrow').hide();
    });
  });

  $('.item').click(function(event){
    var Self = this;
    if(!$(this).find('.drop-content, .drop-arrow').is(':visible')){
      if($(this).is('#cart')){
        $('#cart .drop-body').load('/Cart', function(){
	  $(Self).trigger('mouseenter').find('.cart-total').stop(true, true);
	  $('.drop-content', Self).off('click').on('click', function(event){event.stopPropagation();});
	  $(Self).addClass('active').find('.drop-content, .drop-arrow').show().find('.drop-body').hide().slideDown(200);
      });
      } 
      else if($(this).is('#login')) {
        $('.drop-content', Self).off('click').on('click', function(event){event.stopPropagation();});
	$(this).addClass('active').find('.drop-content, .drop-arrow').show().find('.drop-body').hide().slideDown(100);	
      } 
      else {
        $(this).addClass('active').find('.drop-content, .drop-arrow').show().find('.drop-body').hide().slideDown(200);	
      }
    } 
    else {
      $(this).find('.drop-content .drop-body').slideUp(100, function(){
      $(this).parents('.drop-content').parent().removeClass('active').trigger('mouseleave').find('.drop-content, .drop-arrow').hide();	
      });	
    }
    $('.drop-content .drop-body').not($(this).find('.drop-content .drop-body')).slideUp(100, function(){
      $(this).parents('.drop-content').parent().removeClass('active').trigger('mouseleave').find('.drop-content, .drop-arrow').hide();	
    });	
    $('#search-results-body').slideUp(200, function(){
      $(this).parents('.search-results-content').hide();
    });
    event.stopPropagation();	
  });
	
	
	
  $('#menu .item').mouseenter(function(){
    $(this).find('.light div').stop().animate({'background-color': $('.theme-colors-4').css('backgroundColor')}, 300);
  }).mouseleave(function(){
  if($(this).hasClass('active')) return;
	
  $(this).find('.light div').stop().animate({'background-color': $('.theme-colors-3').css('backgroundColor')}, 400);
	});
	
  $('.nav > ul > li').on('mouseenter mouseleave', function(e){
    var color = (e.type == 'mouseenter') ? '#3333FF' : $('.theme-colors-2').css('backgroundColor');
    var interval = (e.type == 'mouseenter') ? 300 : 400;
    $(this).stop().animate({'background-color': color}, interval);
  });
	
  $('.nav-phone').on('mouseenter mouseleave', function(e){
    var color = (e.type == 'mouseenter') ? '#3333FF' : $('.theme-colors-2').css('backgroundColor');
    var interval = (e.type == 'mouseenter') ? 300 : 400;
    $(this).find('div').stop().animate({'background-color': color}, interval);
  });
	
  $('.grid .struct').live('mouseenter mouseleave', function(e){
    var color1 = (e.type == 'mouseenter') ? '#e5f0ff' : '#fafafa';
    var color2 = (e.type == 'mouseenter') ? '#cce1ff' : '#ebebeb';
    var interval = (e.type == 'mouseenter') ? 300 : 400;
		
    $(this).stop().animate({'background-color': color1}, interval);
    $('.frame', this).stop().animate({'border-color': color2}, interval);
  });
	
	
  $('#cart').mouseenter(function(){
    $(this).find('.cart-total').stop().animate({'background-color': '#66CCFF'}, 300);
  }).mouseleave(function(){
    if($(this).hasClass('active')) return;
    $(this).find('.cart-total').stop().animate({'background-color': '#3333FF'}, 400);
    /*$(this).find('.cart-total').stop().animate({'background-color': $('.theme-colors-2').css('backgroundColor')}, 400);*/
  });
	
  $('.nav > ul > li').hover(function(event){
    var maxWidth = $(this).parents('.nav').outerWidth();
    var menuLeft = $(this).position().left;
	
    $(this).find('.sub-content').css({
      left: -menuLeft,
      width: maxWidth
    });
		
    var menuWidth = $(this).find('.sub-content, .sub-arrow').show().find('.sub-body').css({left: 0}).outerWidth() + 1;

    if(maxWidth <= menuWidth){
      $(this).find('.sub-body').css({left: 0});	
    } else {
      $(this).find('.sub-body').css({
      left: (menuLeft + menuWidth > maxWidth) ? maxWidth - menuWidth : menuLeft
      });	
    }
  }, function(){
    $(this).find('.sub-content, .sub-arrow').hide();	
  });

  $('footer .column h3').click(function(){
    var Parent = $(this).parent();
	
    Parent.find('ul').slideToggle(300, function(){
      Parent.toggleClass('active').find('ul').removeAttr('style');
    });
  });
	
  $('.button').live('mouseenter', function(){
    $(this).stop().animate({'background-color': $('.theme-colors-2').css('backgroundColor')}, 300);
  }).live('mouseleave', function(){
    if($(this).hasClass('button-alt')){
      $(this).stop().animate({'background-color': $('.theme-colors-1').css('backgroundColor')}, 300);		
    } else {
      $(this).stop().animate({'background-color': '#fff'}, 300);	
    }
  });
	
  $('.social a').hover(function(){
    $(this).stop().animate({'opacity': '0.5'}, 300);		
  }, function(){
    $(this).stop().animate({'opacity': '1'}, 300);		
  });
	
  /* Search */
  $('#search .button-search').bind('click', Search);
  $('#search input[name=\'filter_name\']').bind('keydown', function(e){
    if(e.keyCode == 13)
      Search();
  });
	
  $('.messages > div').each(function(i, e){
    Notification($(e).html(), $(e).attr('class'));	
  });
});

function Search(){
        Notification['Search is being implemented', 'success'];
	/*url = $('base').attr('href') + 'product/search';
	var filter_name = $('input[name=\'filter_name\']').val();
	
	if(filter_name && !$('input[name=\'filter_name\']').hasClass('default'))
		url += '&filter_name=' + encodeURIComponent(filter_name);
	
	location = url;	*/
}

function getURLVar(urlVarName) {
	var urlHalves = String(document.location).toLowerCase().split('?');
	var urlVarValue = '';
	
	if (urlHalves[1]) {
		var urlVars = urlHalves[1].split('&');

		for (var i = 0; i <= (urlVars.length); i++) {
			if (urlVars[i]) {
				var urlVarPair = urlVars[i].split('=');
				
				if (urlVarPair[0] && urlVarPair[0] == urlVarName.toLowerCase()) {
					urlVarValue = urlVarPair[1];
				}
			}
		}
	}
	
	return urlVarValue;
} 

function getTotal() {
  $.post("../Cart", {'action':'total'}, function(data) {
    $('.cart-total-text').text(numberWithCommas(data.total)); 
    $('.cart-count').text(data.items); 
  }, "json");
}
function getTotal1() {
  $.post("../../Cart", {'action':'total'}, function(data) {
    $('.cart-total-text').text(numberWithCommas(data.total)); 
    $('.cart-count').text(data.items); 
  }, "json");
}

function numberWithCommas(xx) {
  return xx.toString().replace(/\B(?=(\d{3})+(?!\d))/g,",");
} 

function addToCart(product_id, quantity) {
	quantity = typeof(quantity) != 'undefined' ? quantity : 1;

	$.ajax({
		url: '/Cart',
		type: 'post',
		data: 'action=add' + '&product_id=' + product_id + '&quantity=' + quantity,
		dataType: 'json',
		success: function(json) {
		  if (json['redirect']) {
		    location = json['redirect'];
		  }
		  
		  if (json['success']) {
		    Notification(json['success'], 'success');
		    
		    var matches = /^(.*) - ([^ ]*)$/i.exec(json['total']);
		    
		    $('#cart .cart-count').text(matches[1]);
		    $('#cart .cart-total-text').text(numberWithCommas(matches[2]));
		  }	
		}
	});
}
function addToWishList(product_id) {
	$.ajax({
		url: 'index.php?route=account/wishlist/add',
		type: 'post',
		data: 'product_id=' + product_id,
		dataType: 'json',
		success: function(json) {		
			if (json['success']) {
				Notification(json['success'], 'success');

				$('#wishlist-total').html(json['total']);
			}	
		}
	});
}

function addToCompare(product_id) { 
	$.ajax({
		url: 'index.php?route=product/compare/add',
		type: 'post',
		data: 'product_id=' + product_id,
		dataType: 'json',
		success: function(json) {	
			if (json['success']) {
				Notification(json['success'], 'success');
				
				$('#compare-total').html(json['total']);
			}	
		}
	});
}
