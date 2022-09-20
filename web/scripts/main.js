$(document).ready(function() {
    $("#update").on('click', function() {
        var adv = $("#adv").val();
        var bg = $("#bg").val();
        var state = $("#state").val();
        var apitoken = $("#apitoken").val();
        $.ajax({
            url: "/update",
            method: "GET",
            contentType: "application/x-www-form-urlencoded",
            data: {
                bg: bg,
                state: state,
                apitoken: apitoken,
                adv: adv,
            },
            success: function(data) {
                $("#response").html(data);
            },
        });
    });
});

$(document).ready(function() {
    $("#adminupdate").on('click', function() {
        var adv = $("#adminadv").val();
        var state = $("#adminstate").val();
        var apitoken = $("#admintoken").val();
		var playerid = $("#playerid").val();
        $.ajax({
            url: "/adminUpdate",
            method: "GET",
            contentType: "application/x-www-form-urlencoded",
            data: {
                state: state,
                apitoken: apitoken,
                adv: adv,
				playerid: playerid,
            },
            success: function(data) {
                $("#response").html(data);
            },
        });
    });
});

$(document).ready(function() {
    $("#delete").on('click', function() {
        var token = $("#apitoken").val();
        $.ajax({
            url: "/delete",
            method: "GET",
            contentType: "application/x-www-form-urlencoded",
            data: {
                token: token,
            },
            success: function(data) {
                $("#response").html(data);
            },
        });
    });
});