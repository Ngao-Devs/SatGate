<?php

$finder = (new PhpCsFixer\Finder())
    ->in([__DIR__ . '/services/wordpress-plugin']);

return (new PhpCsFixer\Config())
    ->setRules([
        '@PSR12' => true,
    ])
    ->setFinder($finder);
