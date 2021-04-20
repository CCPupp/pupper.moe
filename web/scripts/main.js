$(document).ready(function() {
    $("#update").on('click', function() {
        var bg = $("#bg").val();
        var mode = $("#mode").val();
        var state = $("#state").val();
        var id = $("#userid").val();
        $.ajax({
            url: "/update",
            method: "GET",
            contentType: "application/x-www-form-urlencoded",
            data: {
                bg: bg,
                mode: mode,
                state: state,
                id: id,
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