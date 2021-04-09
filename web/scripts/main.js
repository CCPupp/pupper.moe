$(document).ready(function() {
    $("#submit").on('click', function() {
        var state = $("#state").val();
        $.ajax({
            url: "/submit",
            method: "GET",
            contentType: "application/x-www-form-urlencoded",
            data: {
                state: state
            },
            success: function(data) {
                $("#response").html(data);
            },
        });
    });
});