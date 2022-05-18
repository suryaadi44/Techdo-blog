$(function () {
    $(document).scroll(function () {
      var $nav = $(".navbar");
      $nav.toggleClass('scrolled', $(this).scrollTop() > $nav.height());
      $(".navbar__create-post-btn").toggleClass('scrolled', $(this).scrollTop() > $nav.height());
      $(".navbar__create-acc-btn").toggleClass('scrolled', $(this).scrollTop() > $nav.height());
      $(".navbar__login-btn").toggleClass('scrolled', $(this).scrollTop() > $nav.height());
      $(".text-logo").toggleClass('scrolled', $(this).scrollTop() > $nav.height());
      $(".half-logo").toggleClass('gradient', $(this).scrollTop() < $nav.height());
      $(".dark-logo").toggleClass('hide', $(this).scrollTop() >= $nav.height());
      $(".white-logo").toggleClass('hide', $(this).scrollTop() < $nav.height());
    });
});