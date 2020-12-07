const baseKey = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/='
let base64 = {
  encode: function ($input) {
    if (!$input) {
      return false
    }
    // default comment
    // $input = UTF8.encode($input);
    let $output = ''
    let $chr1; let $chr2; let $chr3
    let $enc1; let $enc2; let $enc3; let $enc4
    let $i = 0
    do {
      $chr1 = $input.charCodeAt($i++)
      $chr2 = $input.charCodeAt($i++)
      $chr3 = $input.charCodeAt($i++)
      $enc1 = $chr1 >> 2
      $enc2 = (($chr1 & 3) << 4) | ($chr2 >> 4)
      $enc3 = (($chr2 & 15) << 2) | ($chr3 >> 6)
      $enc4 = $chr3 & 63
      if (isNaN($chr2)) $enc3 = $enc4 = 64
      else if (isNaN($chr3)) $enc4 = 64
      $output += baseKey.charAt($enc1) + baseKey.charAt($enc2) + baseKey.charAt($enc3) + baseKey.charAt($enc4)
    } while ($i < $input.length)
    return $output
  },
  decode: function ($input) {
    if (!$input) return false
    $input = $input.replace(/[^A-Za-z0-9+/=]/g, '')
    let $output = ''
    let $enc1; let $enc2; let $enc3; let $enc4
    let $i = 0
    do {
      $enc1 = baseKey.indexOf($input.charAt($i++))
      $enc2 = baseKey.indexOf($input.charAt($i++))
      $enc3 = baseKey.indexOf($input.charAt($i++))
      $enc4 = baseKey.indexOf($input.charAt($i++))
      $output += String.fromCharCode(($enc1 << 2) | ($enc2 >> 4))
      if ($enc3 !== 64) $output += String.fromCharCode((($enc2 & 15) << 4) | ($enc3 >> 2))
      if ($enc4 !== 64) $output += String.fromCharCode((($enc3 & 3) << 6) | $enc4)
    } while ($i < $input.length)
    return $output // UTF8.decode($output);
  }
}

export default base64.encode
