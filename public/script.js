$(document).ready(function () {
    $.ajax({
        url: 'https://test.oma.rs/packs',
        type: 'GET',
        success: function (response) {
            $("#packs").text(JSON.stringify(response.packs).replaceAll(",", ", "));
        }
    });
    $("#order-form").submit(function (event) {
        event.preventDefault();

        var items = $("#order-input").val()
        if (isNaN(items)) {
            alert("Invalid input. Please enter number of items.");
            return;
        }

        $.ajax({
            url: `https://test.oma.rs/order/${items}`,
            type: 'GET',
            contentType: 'application/json;charset=UTF-8',
            success: function (response) {
                alert(response.result);
            }
        });
    });
    $("#write-form").submit(function (event) {
        event.preventDefault();

        var arr = $("#write-input").val().split(",");
        if (!validateInputArray(arr)) {
            alert("Invalid input. Please enter numbers separated by comma.");
            return;
        }

        // Convert array of strings to array of numbers
        arr = arr.map(Number);

        // Sort the array in ascending order
        arr.sort(function (a, b) {
            return a - b;
        });

        // Remove duplicate values from the array
        arr = arr.filter(function (item, pos) {
            return arr.indexOf(item) === pos;
        });

        $.ajax({
            url: 'https://test.oma.rs/packs',
            type: 'POST',
            data: JSON.stringify(arr),
            contentType: 'application/json;charset=UTF-8',
            success: function (response) {
                $("#write-input").val('');
                alert(response.message)
                $.ajax({
                    url: 'https://test.oma.rs/packs',
                    type: 'GET',
                    success: function (response) {
                        $("#packs").text(JSON.stringify(response.packs).replaceAll(",", ", "));
                    }
                });
            }
        });
    });
});


function validateInputArray(arr) {
    for (var i = 0; i < arr.length; i++) {
        if (isNaN(arr[i])) {
            return false;
        }
    }
    return true;
}
