BEGIN { FS = "|"; }

BEGINFILE {
    print "Processing " FILENAME "\n";
}

{
    if($2 ~ /user: /) { # process only lines with "user: "
        user = gensub(/user: *|\s+/, "" , "g", $2);
        success = gensub(/success: *|\s+/, "", "g", $3); 
        if(success == "true") {
            success_count[user]++;
        }
    }
}

ENDFILE {
    totalCount = 0;
    totlaUsers = 0;
    for (user in success_count) {
        printf "%-7s made %-4d req\n", user, success_count[user];
        totalCount += success_count[user];
        totlaUsers++;
        delete success_count[user];
    }
    print "-------------------------";
    print "Average requests per user: " totalCount/totlaUsers;
    print "-------------------------\n";
}