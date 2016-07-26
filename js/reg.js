$(document).ready( function() {
// variable to hold request
var request;
// bind to the submit event of our form
$("#userEntry").submit(function(event){
    //console.log("Hooray, it worked!");
    // abort any pending request
    if (request) {
        request.abort();
    }
    // setup some local variables
    var $form = $(this);
    // let's select and cache all the fields
    var $inputs = $form.find("input, select, button, textarea");
    console.log($inputs);
    // serialize the data in the form
    var serializedData = $form.serialize();
 
    //console.log(serializedData);
    // let's disable the inputs for the duration of the ajax request
    $inputs.prop("disabled", true);

    var fname = $('[id$=fname]').val();			
    var lname = $('[id$=lname]').val();			
    var email = $('[id$=email]').val();			
    var pswd = $('[id$=pswd]').val();			
    var pswd1 = $('[id$=pswd1]').val();			
    var street = $('[id$=street]').val();
    var city = $('[id$=city]').val();	
    var state = $('[id$=state]').val();					
    var cntry = $('[id$=cntry]').val();					
    // fire off the request to /UserServlet
    request = $.ajax({
        url: "/account/register",
        type: "post",
        data: {'action':'register', 'fname':fname, 'lname':lname, 'email':email, 'street':street, 'city':city, 'state':state, 'pswd':pswd, 'pswd1':pswd1, 'cntry':cntry}
    });
    /*$.post("UserServlet", {'action':'register','fullname':name,'email':email,'addr':addr,'area':area,'note':note, 'pswd':pswd},
				function(data){console.log("It worked!"); },"text");
    });*/

    // callback handler that will be called on success
    request.done(function (response, textStatus, jqXHR){
        // log a message to the console
        $('#result').empty();
        $('<p>').text(response).appendTo($('#result'));
        $('<p>').text("Thank you for registering. We'll contact you asap.").appendTo($('#result'));
    });

    // callback handler that will be called on failure
    request.fail(function (jqXHR, textStatus, errorThrown){
        // log the error to the console
        console.error(
            "The following error occured: "+
            textStatus, errorThrown
        );
        $('#result').empty();
        $('<p>').text("Sorry email already in use: " + textStatus, errorThrown ).appendTo($('#result'));
    });

    // callback handler that will be called regardless
    // if the request failed or succeeded
    request.always(function () {
        // reenable the inputs
        $inputs.prop("disabled", false);
    });

    // prevent default posting of form
    event.preventDefault();
});

$("#userreg").submit(function(event) {
    if (request) {
        request.abort();
    }
    var rsa = new RSAKey();
    rsa.setPublic("ac954c14b7d84272710167f354ed64818a091239017ad888a1311145395073b1","10001" );
    document.userreg.password.value = rsa.encrypt(document.userreg.password.value);
    document.userreg.confirm.value = rsa.encrypt(document.userreg.confirm.value);
    // setup some local variables
    var $form = $(this);
    // let's select and cache all the fields
    var $inputs = $form.find("input, select, button, textarea");
    console.log($inputs);
    // serialize the data in the form
    var serializedData = $form.serialize();
    // let's disable the inputs for the duration of the ajax request
    $inputs.prop("disabled", true);

    // fire off the request to /UserServlet
    request = $.ajax({
        url: "/account/register",
        type: "post",
        data: serializedData /*{'action':'login', 'email':email,'pswd':pswd, }*/
    });
    // callback handler that will be called on success
    request.done(function (response, textStatus, jqXHR){
      if (response['success']) {
        location = "/account/account";
      }  
      if (response['failure']) {
	Notification(response['failure'], 'failure');
        document.userreg.password.value = "";
        document.userreg.confirm.value = "";
      }
    });
    // callback handler that will be called on failure
    request.fail(function (jqXHR, textStatus, errorThrown){
        // log the error to the console
        console.error(
            "The following error occured: "+
            textStatus, errorThrown
        );
        location = "/account/register";
    });
    // callback handler that will be called regardless
    // if the request failed or succeeded
    request.always(function () {
        // reenable the inputs
        $inputs.prop("disabled", false);
    });
    // prevent default posting of form
    event.preventDefault();
});

$("#userchange").submit(function(event) {
    if (request) {
        request.abort();
    }
    var rsa = new RSAKey();
    rsa.setPublic("ac954c14b7d84272710167f354ed64818a091239017ad888a1311145395073b1","10001" );
    document.userchange.old_password.value = rsa.encrypt(document.userchange.old_password.value);
    document.userchange.password.value = rsa.encrypt(document.userchange.password.value);
    document.userchange.confirm.value = rsa.encrypt(document.userchange.confirm.value);
    // setup some local variables
    var $form = $(this);
    // let's select and cache all the fields
    var $inputs = $form.find("input, select, button, textarea");
    console.log($inputs);
    // serialize the data in the form
    var serializedData = $form.serialize();
    // let's disable the inputs for the duration of the ajax request
    $inputs.prop("disabled", true);

    // fire off the request to /UserServlet
    request = $.ajax({
        url: "/account/edit",
        type: "post",
        data: serializedData /*{'action':'login', 'email':email,'pswd':pswd, }*/
    });
    // callback handler that will be called on success
    request.done(function (response, textStatus, jqXHR){
      if (response['success']) {
	Notification(response['success'], 'success');
        location = "/account/account";
      }  
      if (response['failure']) {
	Notification(response['failure'], 'failure');
        document.userchange.old_password.value = "";
        document.userchange.password.value = "";
        document.userchange.confirm.value = "";
      }
    });
    // callback handler that will be called on failure
    request.fail(function (jqXHR, textStatus, errorThrown){
        // log the error to the console
        console.error(
            "The following error occured: "+
            textStatus, errorThrown
        );
        location = "/account/register";
    });
    // callback handler that will be called regardless
    // if the request failed or succeeded
    request.always(function () {
        // reenable the inputs
        $inputs.prop("disabled", false);
    });
    // prevent default posting of form
    event.preventDefault();
});

$("#userforgot").submit(function(event) {
    if (request) {
        request.abort();
    }
    // setup some local variables
    var $form = $(this);
    // let's select and cache all the fields
    var $inputs = $form.find("input, select, button, textarea");
    console.log($inputs);
    // serialize the data in the form
    var serializedData = $form.serialize();
    // let's disable the inputs for the duration of the ajax request
    $inputs.prop("disabled", true);

    // fire off the request to /UserServlet
    request = $.ajax({
        url: "/account/forgot",
        type: "post",
        data: serializedData /*{'action':'login', 'email':email,'pswd':pswd, }*/
    });
    // callback handler that will be called on success
    request.done(function (response, textStatus, jqXHR){
      if (response['success']) {
	Notification(response['success'], 'success');
      }  
      if (response['failure']) {
	Notification(response['failure'], 'failure');
      }
    });
    // callback handler that will be called on failure
    request.fail(function (jqXHR, textStatus, errorThrown){
        // log the error to the console
        console.error(
            "The following error occured: "+
            textStatus, errorThrown
        );
        location = "/account/register";
    });
    // callback handler that will be called regardless
    // if the request failed or succeeded
    request.always(function () {
        // reenable the inputs
        $inputs.prop("disabled", false);
    });
    // prevent default posting of form
    event.preventDefault();
});

$("#userlogin").submit(function(event) {
    if (request) {
        request.abort();
    }
    // setup some local variables
    var $form = $(this);
    // let's select and cache all the fields
    var $inputs = $form.find("input, select, button, textarea");
    console.log($inputs);
    // serialize the data in the form
    var serializedData = $form.serialize();
    // let's disable the inputs for the duration of the ajax request
    $inputs.prop("disabled", true);

    var email = $('[id$=email]').val();			
    var pswd = $('[id$=pswd]').val();			
    var rsa = new RSAKey();
    rsa.setPublic("ac954c14b7d84272710167f354ed64818a091239017ad888a1311145395073b1","10001" );
    var res = rsa.encrypt(pswd);
    
    // fire off the request to /UserServlet
    request = $.ajax({
        url: "/account",
        type: "post",
        data: {'action':'login', 'email':email,'pswd':res, }
    });
    // callback handler that will be called on success
    request.done(function (response, textStatus, jqXHR){
        // log a message to the console
        console.log("Successfully logged in");
        location = "/products/front";
    });
    // callback handler that will be called on failure
    request.fail(function (jqXHR, textStatus, errorThrown){
        // log the error to the console
        console.error(
            "The following error occured: "+
            textStatus, errorThrown
        );
        location = "/account/login";
    });
    // callback handler that will be called regardless
    // if the request failed or succeeded
    request.always(function () {
        // reenable the inputs
        $inputs.prop("disabled", false);
    });
    // prevent default posting of form
    event.preventDefault();
});

$("#userlogin1").submit(function(event) {
    if (request) {
        request.abort();
    }
    var rsa = new RSAKey();
    rsa.setPublic("ac954c14b7d84272710167f354ed64818a091239017ad888a1311145395073b1","10001" );
    document.userlogin1.pswd.value = rsa.encrypt(document.userlogin1.pswd.value);
    // setup some local variables
    var $form = $(this);
    // let's select and cache all the fields
    var $inputs = $form.find("input, select, button, textarea");
    console.log($inputs);
    // serialize the data in the form
    var serializedData = $form.serialize();
    // let's disable the inputs for the duration of the ajax request
    $inputs.prop("disabled", true);

    // fire off the request to /UserServlet
    request = $.ajax({
        url: "/account",
        type: "post",
        data: serializedData //{'action':'login', 'email':email,'pswd':pswd, }
    });
    // callback handler that will be called on success
    request.done(function (response, textStatus, jqXHR){
        // log a message to the console
        console.log("Successfully logged in");
        location = "/account/account";
    });
    // callback handler that will be called on failure
    request.fail(function (jqXHR, textStatus, errorThrown){
        // log the error to the console
        console.error(
            "The following error occured: "+
            textStatus, errorThrown
        );
        location = "/account/login";
    });
    // callback handler that will be called regardless
    // if the request failed or succeeded
    request.always(function () {
        // reenable the inputs
        $inputs.prop("disabled", false);
    });
    // prevent default posting of form
    event.preventDefault();
});

$("#userlogout").submit(function(event) {
    if (request) {
        request.abort();
    }
    // setup some local variables
    var $form = $(this);
    // let's select and cache all the fields
    var $inputs = $form.find("input, select, button, textarea");
    console.log($inputs);
    // serialize the data in the form
    var serializedData = $form.serialize();
    // let's disable the inputs for the duration of the ajax request
    $inputs.prop("disabled", true);
    // fire off the request to /UserServlet
    request = $.ajax({
        url: "/account",
        type: "post",
        data: {'action':'logout' }
    });
    // callback handler that will be called on success
    request.done(function (response, textStatus, jqXHR){
        // log a message to the console
        location = "/products/front";
    });
    // callback handler that will be called on failure
    request.fail(function (jqXHR, textStatus, errorThrown){
        // log the error to the console
        console.error(
            "The following error occured: "+
            textStatus, errorThrown
        );
        location = "/products/front";
    });
    // callback handler that will be called regardless
    // if the request failed or succeeded
    request.always(function () {
        // reenable the inputs
        $inputs.prop("disabled", false);
    });
    // prevent default posting of form
    event.preventDefault();
});

$("#product_entry").submit(function(event) {
    if (request) {
        request.abort();
    }
    // setup some local variables
    var $form = $(this);
    // let's select and cache all the fields
    var $inputs = $form.find("input, select, button, textarea");
    // serialize the data in the form
    var serializedData = $form.serialize();
    // let's disable the inputs for the duration of the ajax request
    $inputs.prop("disabled", true);
    // fire off the request to /UserServlet
    request = $.ajax({
      url: "/goods/entry",
      type: "post",
      data: serializedData /*{'action':'login', 'email':email,'pswd':pswd, }*/
    });
    // callback handler that will be called on success
    request.done(function (response, textStatus, jqXHR){
      if (response['success']) {
	Notification(response['success'], 'success');
      }  
      if (response['failure']) {
	Notification(response['failure'], 'failure');
      }
    });
    // callback handler that will be called on failure
    request.fail(function (jqXHR, textStatus, errorThrown){
        // log the error to the console
        console.error(
            "The following error occured: "+
            textStatus, errorThrown
        );
        location = "/account/register";
    });
    // callback handler that will be called regardless
    // if the request failed or succeeded
    request.always(function () {
        // reenable the inputs
        $inputs.prop("disabled", false);
    });
    // prevent default posting of form
    event.preventDefault();
});

$("#product_edit").submit(function(event) {
    if (request) {
        request.abort();
    }
    // setup some local variables
    var $form = $(this);
    // let's select and cache all the fields
    var $inputs = $form.find("input, select, button, textarea");
    // serialize the data in the form
    var serializedData = $form.serialize();
    // let's disable the inputs for the duration of the ajax request
    $inputs.prop("disabled", true);
    // fire off the request 
    request = $.ajax({
      url: "/goods/edit",
      type: "post",
      data: serializedData /*{'action':'login', 'email':email,'pswd':pswd, }*/
    });
    // callback handler that will be called on success
    request.done(function (response, textStatus, jqXHR){
      if (response['success']) {
	Notification(response['success'], 'success');
      }  
      if (response['failure']) {
	Notification(response['failure'], 'failure');
      }
    });
    // callback handler that will be called on failure
    request.fail(function (jqXHR, textStatus, errorThrown){
        // log the error to the console
        console.error(
            "The following error occured: "+
            textStatus, errorThrown
        );
        location = "/account/register";
    });
    // callback handler that will be called regardless
    // if the request failed or succeeded
    request.always(function () {
        // reenable the inputs
        $inputs.prop("disabled", false);
    });
    // prevent default posting of form
    event.preventDefault();
});

$("#product_delete").submit(function(event) {
    if (request) {
        request.abort();
    }
    // setup some local variables
    var $form = $(this);
    // let's select and cache all the fields
    var $inputs = $form.find("input, select, button, textarea");
    // serialize the data in the form
    var serializedData = $form.serialize();
    // let's disable the inputs for the duration of the ajax request
    $inputs.prop("disabled", true);
    // fire off the request 
    request = $.ajax({
      url: "/goods/delete",
      type: "post",
      data: serializedData /*{'action':'login', 'email':email,'pswd':pswd, }*/
    });
    // callback handler that will be called on success
    request.done(function (response, textStatus, jqXHR){
      if (response['success']) {
	Notification(response['success'], 'success');
      }  
      if (response['failure']) {
	Notification(response['failure'], 'failure');
      }
    });
    // callback handler that will be called on failure
    request.fail(function (jqXHR, textStatus, errorThrown){
        // log the error to the console
        console.error(
            "The following error occured: "+
            textStatus, errorThrown
        );
        location = "/account/register";
    });
    // callback handler that will be called regardless
    // if the request failed or succeeded
    request.always(function () {
        // reenable the inputs
        $inputs.prop("disabled", false);
    });
    // prevent default posting of form
    event.preventDefault();
});

$("#goods_entry").submit(function(event) {
    if (request) {
        request.abort();
    }
    // setup some local variables
    var $form = $(this);
    // let's select and cache all the fields
    var $inputs = $form.find("input, select, button, textarea");
    // serialize the data in the form
    var serializedData = $form.serialize();
    // let's disable the inputs for the duration of the ajax request
    $inputs.prop("disabled", true);
    // fire off the request to /UserServlet
    request = $.ajax({
      url: "/goods/entry/new",
      type: "post",
      data: serializedData /*{'action':'login', 'email':email,'pswd':pswd, }*/
    });
    // callback handler that will be called on success
    request.done(function (response, textStatus, jqXHR){
      if (response['success']) {
	Notification(response['success'], 'success');
      }  
      if (response['failure']) {
	Notification(response['failure'], 'failure');
      }
    });
    // callback handler that will be called on failure
    request.fail(function (jqXHR, textStatus, errorThrown){
        // log the error to the console
        console.error(
            "The following error occured: "+
            textStatus, errorThrown
        );
        location = "/account/register";
    });
    // callback handler that will be called regardless
    // if the request failed or succeeded
    request.always(function () {
        // reenable the inputs
        $inputs.prop("disabled", false);
    });
    // prevent default posting of form
    event.preventDefault();
});

$("#goods_entry_new").submit(function(event) {
    if (request) {
        request.abort();
    }
    // setup some local variables
    var $form = $(this);
    // let's select and cache all the fields
    var $inputs = $form.find("input, select, button, textarea");
    // serialize the data in the form
    /*var a = $form.serializeArray();
    var o = {};
    $.each(a, function() {
        if(o[this.name] != undefined) {
            if(!o[this.name].push) {
                o[this.name] = [o[this.name]];
            }
            o[this.name].push(this.value || '');
        }
        else {
            o[this.name] = this.value || '';
        }
    });*/
    var a = $form.serializeJSON();
    // let's disable the inputs for the duration of the ajax request
    $inputs.prop("disabled", true);
    // fire off the request
    request = $.ajax({
      url: "/goods/entry/new1",
      type: "post",
      data: a
      //data: JSON.stringify(o) /*{'action':'login', 'email':email,'pswd':pswd, }*/
    });
    // callback handler that will be called on success
    request.done(function (response, textStatus, jqXHR){
      if (response['success']) {
	Notification(response['success'], 'success');
      }  
      if (response['failure']) {
	Notification(response['failure'], 'failure');
      }
    });
    // callback handler that will be called on failure
    request.fail(function (jqXHR, textStatus, errorThrown){
        // log the error to the console
        console.error(
            "The following error occured: "+
            textStatus, errorThrown
        );
        location = "/account/register";
    });
    // callback handler that will be called regardless
    // if the request failed or succeeded
    request.always(function () {
        // reenable the inputs
        $inputs.prop("disabled", false);
    });
    // prevent default posting of form
    event.preventDefault();
});
});

function logout() {
    // fire off the request to /UserServlet
    request = $.ajax({
        url: "/account",
        type: "post",
        data: {'action':'logout' }
    });
    // callback handler that will be called on success
    request.done(function (response, textStatus, jqXHR){
        // log a message to the console
        console.log("Successfully logged out");
        location = "/account/login";
    });
    // callback handler that will be called on failure
    request.fail(function (jqXHR, textStatus, errorThrown){
        // log the error to the console
        console.error(
            "The following error occured: "+
            textStatus, errorThrown
        );
        location = "/account/login";
    });
    // callback handler that will be called regardless
    // if the request failed or succeeded
    request.always(function () {
        // reenable the inputs
        //$inputs.prop("disabled", false);
    });
}

