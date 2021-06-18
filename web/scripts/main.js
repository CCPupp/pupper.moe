$(document).ready(function() {
    $("#update").on('click', function() {
        var adv = $("#adv").val();
        var bg = $("#bg").val();
        var state = $("#state").val();
        var id = $("#userid").val();
        $.ajax({
            url: "/update",
            method: "GET",
            contentType: "application/x-www-form-urlencoded",
            data: {
                bg: bg,
                state: state,
                id: id,
                adv: adv,
            },
            success: function(data) {
                $("#response").html(data);
            },
        });
    });
});

$(document).ready(function() {
    $("#submitPlayer").on('click', function() {
        var state = $("#state").val();
        var id = $("#id").val();
        $.ajax({
            url: "/submitPlayer",
            method: "GET",
            contentType: "application/x-www-form-urlencoded",
            data: {
                state: state,
                id: id
            },
            success: function(data) {
                $("#response").html(data);
            },
        });
    });
});

$(document).ready(function() {
    $("#adminupdate").on('click', function() {
        var id = $("#userid").val();
        $.ajax({
            url: "/adminUpdate",
            method: "GET",
            contentType: "application/x-www-form-urlencoded",
            data: {
                id: id
            },
            success: function(data) {
                $("#testing-panel").html(data);
            },
        });
    });
});