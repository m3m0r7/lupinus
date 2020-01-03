<?php
require __DIR__ . '/vendor/autoload.php';

$env = Dotenv\Dotenv::create(__DIR__);
$env->load();

$based_directory = $_ENV['LILY_CERTIFICATE_BASED_DIRECTORY'] ?? null;
$target = $_ENV['LILY_NGINX_TARGET'] ?? null;

exec('whoami', $output);

if (($output[0] ?? null) !== 'root') {
    echo "You need to login with root user.\n";
    exit(1);
}

if ($based_directory === null || $target === null) {
    echo "Please set envs LILY_CERTIFICATE_BASED_DIRECTORY and LILY_NGINX_TARGET.\n";
    exit(1);
}

$maps = [
    'fullchain.pem' => 'private.pem',
    'privkey.pem' => 'private.key',
];

foreach (glob($based_directory . '/{fullchain,privkey}.pem', GLOB_BRACE) as $file) {
    $name = basename($file);
    echo "Copying $name\n";
    file_put_contents(
        $target . '/' . $maps[$name],
        file_get_contents($file)
    );
}

echo "Start to rebuild nginx.\n";
system("cd $target && docker build . --no-cache");
echo "Finish to rebuild.\n";
echo "Please restart server with `sudo reboot`.\n";
