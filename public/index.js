$(document).ready(function(){
  $("li").click(function(){
    var id = $(this).attr('id');
    $.ajax({
      url: '/'+id,
      type: 'DELETE',
      success: function(res) {
        $('#'+id).remove();
      }
    })
  })
})
